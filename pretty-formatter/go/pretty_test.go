package pretty

import (
	"testing"
	"bytes"
	gio "github.com/gogo/protobuf/io"
	"github.com/stretchr/testify/require"
	//"github.com/fatih/color"
	"github.com/cucumber/cucumber-messages-go"
	"github.com/cucumber/gherkin-go"
)

func TestPrintsFeatureHeaderWithComments(t *testing.T) {
	src := `# Hello
Feature: Hello
`

	stdout := &bytes.Buffer{}
	ProcessMessages(messageReader(t, src), stdout, false)

	require.EqualValues(t,
		src,
		stdout.String())
}

func TestPrintsFeatureAndScenarioHeadersWithComments(t *testing.T) {
	src := `# Hello
Feature: Hello

  # World
  Scenario: World
`

	stdout := &bytes.Buffer{}
	ProcessMessages(messageReader(t, src), stdout, false)

	require.EqualValues(t,
		src,
		stdout.String())
}

func TestPrintsAllTheThings(t *testing.T) {
	src := `# A
Feature: A

  # B
  Background: B
    Given b
      | text | number |
      | a    |     10 |
      | bb   |    100 |
      | ccc  |   1000 |

  # C
  Scenario: C
    Given c
      """
      x
       y
        z
      """
    And <c1>
    Then <c2>

    # CE
    @ce
    Examples: CE
      | c1  |   c2 |
      | a   |   10 |
      | bb  |  100 |
      | ccc | 1000 |

  # D
  Rule: D

    # E
    Background: E
      Given e

    # F
    @f @F
    Scenario: F
      Given f
`

	stdout := &bytes.Buffer{}
	ProcessMessages(messageReader(t, src), stdout, false)

	require.EqualValues(t,
		src,
		stdout.String())
}

func TestPrintsInResultsMode(t *testing.T) {
	src := `Feature: A

  Scenario: B
    passed Given C
    failed When D
    undefined Then E
`

	stdout := &bytes.Buffer{}
	prettyStdin := messageReader(t, src)
	prettyStdinWriter := gio.NewDelimitedWriter(prettyStdin)
	prettyStdinWriter.WriteMsg(newTestStepFinished(messages.Status_PASSED, 3, 0))
	prettyStdinWriter.WriteMsg(newTestStepFinished(messages.Status_FAILED, 3, 1))
	prettyStdinWriter.WriteMsg(newTestStepFinished(messages.Status_UNDEFINED, 3, 2))

	ProcessMessages(prettyStdin, stdout, true)

	require.EqualValues(t,
		src,
		stdout.String())
}

func messageReader(t *testing.T, src string) *bytes.Buffer {
	source := &messages.Source{
		Uri:  "features/test.feature",
		Data: src,
		Media: &messages.Media{
			Encoding:    "UTF-8",
			ContentType: "text/x.cucumber.gherkin+plain",
		},
	}
	sources := &bytes.Buffer{}
	sourcesWriter := gio.NewDelimitedWriter(sources)
	sourcesWriter.WriteMsg(source)
	wrappers, err := gherkin.GherkinMessages(
		nil,
		sources,
		"en",
		true,
		true,
		true,
	)
	require.NoError(t, err)
	prettyStdin := &bytes.Buffer{}
	prettyStdinWriter := gio.NewDelimitedWriter(prettyStdin)
	for _, wrapper := range wrappers {
		prettyStdinWriter.WriteMsg(&wrapper)
	}
	return prettyStdin
}

func newTestStepFinished(status messages.Status, line uint32, index uint32) *messages.Wrapper {
	return &messages.Wrapper{
		Message: &messages.Wrapper_TestStepFinished{
			TestStepFinished: &messages.TestStepFinished{
				TestCase: &messages.TestCase{
					SourceLine: &messages.SourceLine{
						Uri:  "features/test.feature",
						Line: line,
					},
				},
				TestResult: &messages.TestResult{
					Status: status,
				},
				Index: index,
			},
		},
	}
}