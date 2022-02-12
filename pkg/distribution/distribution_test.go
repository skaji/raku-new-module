package distribution

import (
	"encoding/json"
	"testing"
)

func TestDistribution(t *testing.T) {
	j := `
	{
      "id": "zef:lizmat/Identity::Utils/0.0.9",
      "content_text": "Provide utility functions related to distribution identities",
      "title": "Identity::Utils 0.0.9",
      "authors": [
        {
          "url": "https://raku.land/zef:lizmat",
          "name": "Elizabeth Mattijsen",
          "avatar": "https://www.gravatar.com/avatar/d40db3cabae1b579841c2e60f099529c?d=identicon"
        }
      ],
      "date_published": "2022-02-10T17:47:35Z",
      "url": "https://raku.land/zef:lizmat/Identity::Utils"
    }
	`
	var d *Distribution
	if err := json.Unmarshal([]byte(j), &d); err != nil {
		t.Fatal(err)
	}
	if d.Published.Unix() != 1644515255 {
		t.Fatal()
	}
	if d.Auth() != "zef:lizmat" {
		t.Fatal()
	}
}
