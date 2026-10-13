package uiexamples

import (
	"html"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFormsVisual7SourceSpecimens(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatal(err)
	}
	cases := map[string][]string{
		"currencyinput-percentinput-unitinput": {`value="1249.50"`, `value="8"`, `value="16"`, `data-ely-unit-option="em"`, `is-width-200`, `is-width-160`, `is-width-180`},
		"emailinput-urlinput":                  {`type="email"`, `type="url"`, `placeholder="name@example.com"`, `placeholder="https://"`, `is-width-280`},
		"phoneinput-maskedinput":               {`data-ely-phone-country`, `🇺🇸 +1`, `data-ely-mask="(999) 999-9999"`, `data-ely-mask="aa-9999"`, `masked: aa-9999`},
		"radio-radiogroup":                     {`Enterprise`, `aria-disabled="true"`, `density, in a row`, `data-preview-stack`},
		"checkboxcard-radiocard":               {`Personal`, `Team`, `Just you. Notes stay on this Mac.`, `Shared with everyone you invite.`, `Back up nightly`, `type="radio"`, `type="checkbox"`},
		"rating":                               {`3 of 5`, `read only`, `data-ely-rate`},
		"knob-dial-stepper":                    {`gain 65`, `pan 0`, `guests`, `data-ely-step="up"`},
	}
	for slug, wants := range cases {
		t.Run(slug, func(t *testing.T) {
			w := httptest.NewRecorder()
			s.ServeHTTP(w, httptest.NewRequest("GET", "/components/forms/"+slug, nil))
			if w.Code != 200 {
				t.Fatalf("status %d: %s", w.Code, w.Body.String())
			}
			for _, want := range wants {
				if !strings.Contains(html.UnescapeString(w.Body.String()), want) {
					t.Errorf("missing %s", want)
				}
			}
		})
	}
}
