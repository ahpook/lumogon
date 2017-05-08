package analytics

import (
	"testing"

	gock "gopkg.in/h2non/gock.v1"
)

// GockBootstrap returns a UserSession ready for use in testing the GA service.
// Note: This endpoint ALWAYS returns 200 OK, and sends back a small 1px gif image.
// This means we must rely on whether our hit is put together properly.
// Validate at https://ga-dev-tools.appspot.com/hit-builder/
func GockBootstrap() *UserSession {
	u := *NewUserSession()
	gock.New("https://www.google-analytics.com").
		Post("/collect").
		Reply(200).
		SetHeader("Content-Type", "image/gif").
		BodyString("GIF89a�����,D;")
	gock.InterceptClient(u.HTTPClient)
	return &u
}

func GockTeardown(t *testing.T) {
	if gock.IsDone() != true {
		t.Errorf("Pending mock requests are still in queue")
	}
}

func TestScreenViewPostMeasurement(t *testing.T) {
	defer gock.Off()

	u := *GockBootstrap()
	u.Type = "screenview"
	u.ScreenName = "testscreen"
	res, err := u.PostMeasurement()
	if err != nil {
		t.Errorf("HTTP Client did not return 200. How did you get here?")
	}

	analyticsValues := res.Request.URL.Query()
	if len(analyticsValues) != 7 {
		t.Errorf("Does not match expected Query length of 7, returned %d", len(analyticsValues))
		t.Error(analyticsValues)
	}

	GockTeardown(t)
}

func TestEventPostMeasurement(t *testing.T) {
	defer gock.Off()

	u := *GockBootstrap()
	u.Type = "event"
	u.Action = "testaction"
	u.Category = "testcategory"

	res, err := u.PostMeasurement()
	if err != nil {
		t.Errorf("HTTP Client did not return 200. How did you get here?")
	}

	analyticsValues := res.Request.URL.Query()
	if len(analyticsValues) != 8 {
		t.Errorf("Does not match expected Query length of 8, returned %d", len(analyticsValues))
		t.Error(analyticsValues)
	}

	GockTeardown(t)
}

func TestDisableTransmitPostMeasurement(t *testing.T) {
	defer gock.Off()

	t.Skipf("Skipping temporarily")

}
