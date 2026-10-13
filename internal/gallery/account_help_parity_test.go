package uiexamples

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAccountHelpSourceStories(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatal(err)
	}
	cases := map[string][]string{
		"settings/settingslayout-settingssection-settingsrow-settingssearch": {"Appearance", "Keyboard", "Network", "Privacy", "Notifications", "Startup", "Storage", "Advanced", "Buttons, links and what is chosen.", `value="14"`, "data-suite-proxy", "proxy_host", "proxy_port", "mentions_email", "replies_push", "digest_app", "data-suite-reset-defaults", `data-ely-action="clear-storage"`},
		"account/loginform-oauthbuttons":                                     {"Sign in to Ely", "Continue with GitHub", "Continue with Google", "Continue with Apple", "Keep me signed in", "Forgot password?"},
		"account/signupform":                                                 {"Create an account", "At least 8 characters", "A letter", "A number", "I agree to the terms of service"},
		"account/magiclinkform-forgotpassword":                               {"Sign in by email", "Reset your password", "No password needed.", "Send reset link"},
		"account/twofactorinput":                                             {"Two-step sign-in", "Enter the six-digit code", "Use a recovery code"},
		"account/accountswitcher-workspaceswitcher-orgswitcher-usermenu":     {"Atlas", "Pro · 12 members", "Field Notes", "Free · 3 members", "ada@atlas.dev", "ada@hey.com", "data-suite-menu-toggle"},
		"account/profilecard-profileeditor":                                  {"Ada Lovelace", "Analyst, Atlas", "London", "A line or two others see.", `data-suite-validate="profile"`},
		"account/sessionlist-apikeymanager":                                  {"MacBook Pro", "iPhone 17", "Office desktop", "ely_live_…4f2a", "ely_live_…91c0", "Sign out of all other devices"},
		"account/subscriptioncard-planbadge-usagequota-paymentmethodform":    {"Renews on 1 October 2026.", "Seats across every workspace you own", "96.5 of 100 GB", "Resets on 1 October.", "Name on card", `data-suite-validate="payment"`},
		"account/billinghistory-upgradeprompt-paywall":                       {"1 Sep 2026", "1 Aug 2026", "1 Jul 2026", "12 Jun 2026", "Refunded", "9 of 10 projects used", "Version history is on Pro", "Compare any two"},
		"account/teammembertable-roleselector-invitationlist":                {"Ada Lovelace", "Grace Hopper", "Alan Turing", "Owner", "katherine@nasa.gov", "edsger@tue.nl", "Expired"},
		"onboarding/onboardingwizard-welcometour":                            {"Name your workspace", "Bring your team", "You&#39;re set", "Start over", "data-suite-wizard-done"},
		"onboarding/coachmark-hotspot-featurehighlight-newbadge":             {"Pin a view", "Inbox", "Automations", "Split view", "Show again"},
		"onboarding/checklist-first-run-tasks":                               {"Add your photo", "Invite a teammate", "Connect your calendar", "Learn three shortcuts", "So your team knows who is who."},
		"onboarding/helptooltip-inlinehelp-emptystatewithaction":             {"Retention", "Billing email", "Only admins can change billing.", "No projects yet", "New project"},
		"onboarding/helppanel-docssidebar":                                   {"Getting started", "Invite your team", "Change billing", "Work from the keyboard", "Only admins can change billing.", "data-suite-restore"},
		"onboarding/keyboardshortcutcheatsheet-whatsnewdialog":               {"Open anything", "Every shortcut", "Zoom in", "Sep 25, 2026", "Sep 9, 2026", "Sync resumes after sleep.", "data-suite-shortcuts-open"},
		"onboarding/feedbackwidget-contactsupport":                           {"data-suite-feedback-open", "Choose a topic", "Something broke", "My account", "Something else", "help@example.com", "What happened, and what you expected"},
	}
	for path, wants := range cases {
		t.Run(path, func(t *testing.T) {
			r := httptest.NewRecorder()
			s.ServeHTTP(r, httptest.NewRequest("GET", "/components/"+path, nil))
			if r.Code != 200 {
				t.Fatalf("status %d: %s", r.Code, r.Body.String())
			}
			for _, want := range wants {
				if !strings.Contains(r.Body.String(), want) {
					t.Errorf("missing %q", want)
				}
			}
		})
	}
}
