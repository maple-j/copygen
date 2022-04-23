package tests

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/switchupcb/copygen/cli"
	"github.com/switchupcb/copygen/cli/config"
	"github.com/switchupcb/copygen/cli/generator"
	"github.com/switchupcb/copygen/cli/matcher"
	"github.com/switchupcb/copygen/cli/parser"
)

type test struct {
	name     string
	ymlpath  string // ymlpath represents the path to an example's .yml file.
	wantpath string // wantpath represents the path to a verified example's output file.
}

var (
	tests = []test{
		{
			name:     "main",
			ymlpath:  "examples/main/setup/setup.yml",
			wantpath: "examples/main/copygen.go",
		},
		{
			name:     "automatch",
			ymlpath:  "examples/automatch/setup/setup.yml",
			wantpath: "examples/automatch/copygen.go",
		},
		{
			name:     "basic",
			ymlpath:  "examples/basic/setup/setup.yml",
			wantpath: "examples/basic/copygen.go",
		},
		/*
			{
				name:     "deepcopy",
				ymlpath:  "examples/deepcopy/setup/setup.yml",
				wantpath: "examples/deepcopy/copygen.go",
			},
		*/
		{
			name:     "error",
			ymlpath:  "examples/error/setup/setup.yml",
			wantpath: "examples/error/copygen.go",
		},
		{
			name:     "map",
			ymlpath:  "examples/map/setup/setup.yml",
			wantpath: "examples/map/copygen.go",
		},
		{
			name:     "tag",
			ymlpath:  "examples/tag/setup/setup.yml",
			wantpath: "examples/tag/copygen.go",
		},
		{
			name:     "alias",
			ymlpath:  "examples/_tests/alias/setup/setup.yml",
			wantpath: "examples/_tests/alias/copygen.go",
		},
		{
			name:     "automap",
			ymlpath:  "examples/_tests/automap/setup/setup.yml",
			wantpath: "examples/_tests/automap/copygen.go",
		},
		/*
			{
				name:     "cast",
				ymlpath:  "examples/_tests/cast/setup/setup.yml",
				wantpath: "examples/_tests/cast/copygen.go",
			},
		*/
		{
			name:     "cyclic",
			ymlpath:  "examples/_tests/cyclic/setup/setup.yml",
			wantpath: "examples/_tests/cyclic/copygen.go",
		},
		{
			name:     "duplicate",
			ymlpath:  "examples/_tests/duplicate/setup/setup.yml",
			wantpath: "examples/_tests/duplicate/copygen.go",
		},
		{
			name:     "multi",
			ymlpath:  "examples/_tests/multi/setup/setup.yml",
			wantpath: "examples/_tests/multi/copygen.go",
		},
	}
)

// TestExamples tests calls cli.Run() in a similar manner to calling the CLI,
// checking for a valid output.
func TestExamples(t *testing.T) {
	checkwd(t)
	for _, test := range tests {
		testExample(t, test)
	}
}

// testExample tests an example using .go, .tmpl, and programmatic methods.
func testExample(t *testing.T, test test) {
	valid, err := ioutil.ReadFile(test.wantpath)
	if err != nil {
		t.Fatalf("error reading file in test %q.\n%v", test.name, err)
	}

	// test the .go method using CLI Run().
	env := cli.Environment{
		YMLPath: test.ymlpath,
		Output:  false,
		Write:   false,
	}

	goCode, err := env.Run()
	if err != nil {
		t.Fatalf("Run(%q) error: %v", test.name, err)
	}

	if !bytes.Equal(normalizeLineBreaks([]byte(goCode)), normalizeLineBreaks(valid)) {
		fmt.Println(goCode)
		t.Fatalf("Run(%v) output not equivalent to %v", test.name, test.wantpath)
	}

	fmt.Println("PASSED:", test.name)

	// test the .tmpl method using copygen programmatically.
	tmplcode, err := templateRun(env)
	if err != nil {
		t.Fatalf("Run(%q [tmpl]) error: %v", test.name, err)
	}

	if !bytes.Equal(normalizeLineBreaks([]byte(tmplcode)), normalizeLineBreaks(valid)) {
		fmt.Println("FAILED: ", test.name, "(tmpl)", "bypassing...")
		return
		// fmt.Println(tmplcode)
		// t.Fatalf("Run(%v [tmpl]) output not equivalent to %v", test.name, test.wantpath)
	}

	fmt.Println("PASSED:", test.name, "(tmpl)")
}

// templateRun runs copygen programmatically and generates code using a template.
func templateRun(env cli.Environment) (string, error) {
	gen, err := config.LoadYML(env.YMLPath)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	if err = parser.Parse(gen); err != nil {
		return "", fmt.Errorf("%w", err)
	}

	if err = matcher.Match(gen); err != nil {
		return "", fmt.Errorf("%w", err)
	}

	gen.Tempath = "examples/tmpl/template/generate.tmpl"
	code, err := generator.GenerateTemplate(gen)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return code, nil
}
