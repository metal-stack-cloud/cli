package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"slices"

	"bou.ke/monkey"
	"github.com/go-openapi/strfmt"
	"github.com/google/go-cmp/cmp"
	apitests "github.com/metal-stack-cloud/api/go/tests"
	"github.com/metal-stack-cloud/cli/cmd/completion"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/metal-stack/metal-lib/pkg/testcommon"
	"github.com/spf13/afero"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"
)

var testTime = time.Date(2022, time.May, 19, 1, 2, 3, 4, time.UTC)

func init() {
	_ = monkey.Patch(time.Now, func() time.Time { return testTime })
}

type Test[R any] struct {
	Name string
	Cmd  func(want R) []string

	ClientMocks *apitests.ClientMockFns
	FsMocks     func(fs afero.Fs, want R)

	DisableMockClient bool // can switch off mock client creation

	WantErr       error
	Want          R       // for json and yaml
	WantTable     *string // for table printer
	WantWideTable *string // for wide table printer
	Template      *string // for template printer
	WantTemplate  *string // for template printer
	WantMarkdown  *string // for markdown printer
}

func (c *Test[R]) TestCmd(t *testing.T) {
	require.NotEmpty(t, c.Name, "test name must not be empty")
	require.NotEmpty(t, c.Cmd, "cmd must not be empty")

	if c.WantErr != nil {
		_, _, conf := c.newMockConfig(t)

		cmd := newRootCmd(conf)
		os.Args = append([]string{config.BinaryName}, c.Cmd(c.Want)...)

		err := cmd.Execute()
		if diff := cmp.Diff(c.WantErr, err, testcommon.IgnoreUnexported(), testcommon.ErrorStringComparer()); diff != "" {
			t.Errorf("error diff (+got -want):\n %s", diff)
		}
	}

	for _, format := range outputFormats(c) {
		format := format
		t.Run(fmt.Sprintf("%v", format.Args()), func(t *testing.T) {
			_, out, conf := c.newMockConfig(t)

			cmd := newRootCmd(conf)
			os.Args = append([]string{config.BinaryName}, c.Cmd(c.Want)...)
			os.Args = append(os.Args, format.Args()...)

			err := cmd.Execute()
			require.NoError(t, err)

			format.Validate(t, out.Bytes())
		})
	}
}

func (c *Test[R]) newMockConfig(t *testing.T) (any, *bytes.Buffer, *config.Config) {
	mock := apitests.New(t)

	fs := afero.NewMemMapFs()
	if c.FsMocks != nil {
		c.FsMocks(fs, c.Want)
	}

	var (
		out    bytes.Buffer
		config = &config.Config{
			Fs:         fs,
			Out:        &out,
			Completion: &completion.Completion{},
			Client:     mock.Client(c.ClientMocks),
		}
	)

	if c.DisableMockClient {
		config.Client = nil
	}

	return nil, &out, config
}

func AssertExhaustiveArgs(t *testing.T, args []string, exclude ...string) {
	assertContainsPrefix := func(ss []string, prefix string) error {
		for _, s := range ss {
			if strings.HasPrefix(s, prefix) {
				return nil
			}
		}
		return fmt.Errorf("not exhaustive: does not contain " + prefix)
	}

	root := newRootCmd(&config.Config{})
	cmd, args, err := root.Find(args)
	require.NoError(t, err)

	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		if slices.Contains(exclude, f.Name) {
			return
		}
		require.NoError(t, assertContainsPrefix(args, "--"+f.Name), "please ensure you all available args are used in order to increase coverage or exclude them explicitly")
	})
}

func MustMarshal(t *testing.T, d any) []byte {
	b, err := json.MarshalIndent(d, "", "    ")
	require.NoError(t, err)
	return b
}

func MustMarshalToMultiYAML[R any](t *testing.T, data []R) []byte {
	var parts []string
	for _, elem := range data {
		parts = append(parts, string(MustMarshal(t, elem)))
	}
	return []byte(strings.Join(parts, "\n---\n"))
}

func MustJsonDeepCopy[O any](t *testing.T, object O) O {
	raw, err := json.Marshal(&object)
	require.NoError(t, err)
	var copy O
	err = json.Unmarshal(raw, &copy)
	require.NoError(t, err)
	return copy
}

func outputFormats[R any](c *Test[R]) []outputFormat[R] {
	var formats []outputFormat[R]

	if !pointer.IsZero(c.Want) {
		formats = append(formats, &jsonOutputFormat[R]{want: c.Want}, &yamlOutputFormat[R]{want: c.Want})
	}

	if c.WantTable != nil {
		formats = append(formats, &tableOutputFormat[R]{table: *c.WantTable})
	}

	if c.WantWideTable != nil {
		formats = append(formats, &wideTableOutputFormat[R]{table: *c.WantWideTable})
	}

	if c.Template != nil && c.WantTemplate != nil {
		formats = append(formats, &templateOutputFormat[R]{template: *c.Template, templateOutput: *c.WantTemplate})
	}

	if c.WantMarkdown != nil {
		formats = append(formats, &markdownOutputFormat[R]{table: *c.WantMarkdown})
	}

	return formats
}

type outputFormat[R any] interface {
	Args() []string
	Validate(t *testing.T, output []byte)
}

type jsonOutputFormat[R any] struct {
	want R
}

func (o *jsonOutputFormat[R]) Args() []string {
	return []string{"-o", "jsonraw"}
}

func StrFmtPtrDateComparer() cmp.Option {
	return cmp.Comparer(func(x, y *strfmt.DateTime) bool {
		if x == nil && y == nil {
			return true
		}
		if x == nil && y != nil {
			return false
		}
		if x != nil && y == nil {
			return false
		}
		return time.Time(*x).Unix() == time.Time(*y).Unix()
	})
}

func (o *jsonOutputFormat[R]) Validate(t *testing.T, output []byte) {
	var got R

	err := json.Unmarshal(output, &got)
	require.NoError(t, err, string(output))

	if diff := cmp.Diff(o.want, got, testcommon.IgnoreUnexported(), testcommon.StrFmtDateComparer()); diff != "" {
		t.Errorf("diff (+got -want):\n %s", diff)
	}
}

type yamlOutputFormat[R any] struct {
	want R
}

func (o *yamlOutputFormat[R]) Args() []string {
	return []string{"-o", "yamlraw"}
}

func (o *yamlOutputFormat[R]) Validate(t *testing.T, output []byte) {
	var got R

	err := yaml.Unmarshal(output, &got)
	require.NoError(t, err)

	if diff := cmp.Diff(o.want, got, testcommon.IgnoreUnexported(), testcommon.StrFmtDateComparer()); diff != "" {
		t.Errorf("diff (+got -want):\n %s", diff)
	}
}

type tableOutputFormat[R any] struct {
	table string
}

func (o *tableOutputFormat[R]) Args() []string {
	return []string{"-o", "table"}
}

func (o *tableOutputFormat[R]) Validate(t *testing.T, output []byte) {
	validateTableRows(t, o.table, string(output))
}

type wideTableOutputFormat[R any] struct {
	table string
}

func (o *wideTableOutputFormat[R]) Args() []string {
	return []string{"-o", "wide"}
}

func (o *wideTableOutputFormat[R]) Validate(t *testing.T, output []byte) {
	validateTableRows(t, o.table, string(output))
}

type templateOutputFormat[R any] struct {
	template       string
	templateOutput string
}

func (o *templateOutputFormat[R]) Args() []string {
	return []string{"-o", "template", "--template", o.template}
}

func (o *templateOutputFormat[R]) Validate(t *testing.T, output []byte) {
	t.Logf("got following template output:\n\n%s\n\nconsider using this for test comparison if it looks correct.", string(output))

	if diff := cmp.Diff(strings.TrimSpace(o.templateOutput), strings.TrimSpace(string(output))); diff != "" {
		t.Errorf("diff (+got -want):\n %s", diff)
	}
}

type markdownOutputFormat[R any] struct {
	table string
}

func (o *markdownOutputFormat[R]) Args() []string {
	return []string{"-o", "markdown"}
}

func (o *markdownOutputFormat[R]) Validate(t *testing.T, output []byte) {
	validateTableRows(t, o.table, string(output))
}

func validateTableRows(t *testing.T, want, got string) {
	trimAll := func(ss []string) []string {
		var res []string
		for _, s := range ss {
			res = append(res, strings.TrimSpace(s))
		}
		return res
	}

	var (
		trimmedWant = strings.TrimSpace(want)
		trimmedGot  = strings.TrimSpace(string(got))

		wantRows = trimAll(strings.Split(trimmedWant, "\n"))
		gotRows  = trimAll(strings.Split(trimmedGot, "\n"))
	)

	t.Logf("got following table output:\n\n%s\n\nconsider using this for test comparison if it looks correct.", trimmedGot)

	t.Log(cmp.Diff(trimmedWant, trimmedGot))

	require.Equal(t, len(wantRows), len(gotRows), "tables have different lengths")

	for i := range wantRows {
		wantFields := trimAll(strings.Split(wantRows[i], " "))
		gotFields := trimAll(strings.Split(gotRows[i], " "))

		require.Equal(t, len(wantFields), len(gotFields), "table fields have different lengths")

		for i := range wantFields {
			assert.Equal(t, wantFields[i], gotFields[i])
		}
	}
}

func commonExcludedFileArgs() []string {
	return []string{"file", "bulk-output", "skip-security-prompts", "timestamps"}
}
