package controller

import (
	"reflect"
	"testing"

	"github.com/whatflix/entity"
)

//const fakeToken = "fake_session_token"

/*func TestRecommendHTTPHandler(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/movies/users", nil)

	h := newHTTPHandler(
		cfg:                    cfg,
		userPreferencesManager: &userPreferencesManagerImpl{},
		creditsManager:         &creditsManagerImpl{},
		moviesManager:          &moviesManagerImpl{},
	)
}*/

/*func TestRecommendHTTPHandler(t *testing.T) {
server := httptest.NewServer(&recommendHTTPHandler{})
defer server.Close()

client := signedInClient(t, server.URL)
_ = client
resp, err := client.Get(server.URL + "/admin")
if err != nil {
	t.Errorf("GET /admin err = %s; want nil", err)
}
defer resp.Body.Close()

if resp.StatusCode != 200 {
	t.Errorf("GET /admin code = %d", resp.StatusCode, 200)
}*/
/*body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutill.ReadAll() err = %s; want nil", err)
	}
	got := string(body)

}*/

/*func signedInClient(t *testing.T, baseURL string) *http.Client {
	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		t.Fatalf("cookiejar.New() err = %s; want nil", err)
	}

	client := &http.Client{
		Jar: jar,
	}

	loginURL := baseURL + "/movies/signin"
	data := url.Values{}
	data.Set("username", "admin")
	data.Set("password", "admin")
	req, err := http.NewRequest(http.MethodPost, loginURL, strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatalf("NewRequest() err = %s; want nil", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("POST /login err= %s; want nil", err)
	}
	t.Log("resp ", resp)
	t.Logf("Cookies %v", client.Jar.Cookies(req.URL))
	return client
}*/

func TestRemoveDuplicate(t *testing.T) {
	var entries []*entity.CreditsData
	entries = append(entries, &entity.CreditsData{
		Title: "Avengers"})
	entries = append(entries, &entity.CreditsData{
		Title: "The Hurt Locker"})
	entries = append(entries, &entity.CreditsData{
		Title: "Avengers"})
	entries = append(entries, &entity.CreditsData{
		Title: "Mr. & Mrs. Smith"})

	want := []string{"Avengers", "The Hurt Locker", "Mr. & Mrs. Smith"}
	got := removeDuplicate(entries)

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected got = %s\n want = %s", got, want)
	}
}
