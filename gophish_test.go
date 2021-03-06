package gophish

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestParseDate(t *testing.T) {
	tests := []struct {
		date string
		want time.Time
	}{
		{
			date: "1980-10-03",
			want: time.Date(1980, time.October, 3, 0, 0, 0, 0, time.UTC),
		}, {
			date: "2000-01-01",
			want: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tc := range tests {
		got, err := ParseDate(tc.date)
		if err != nil {
			t.Errorf("Got unexpected err; %v", err)
		} else if !got.Equal(tc.want) {
			t.Errorf("Got: %s\nWant: %s", got, tc.want)
		}
	}
}

func TestFormatDate(t *testing.T) {
	tests := []struct {
		want string
		date time.Time
	}{
		{
			want: "1980-10-03",
			date: time.Date(1980, time.October, 3, 0, 0, 0, 0, time.UTC),
		}, {
			want: "2000-01-01",
			date: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tc := range tests {
		got := FormatDate(tc.date)
		if got != tc.want {
			t.Errorf("Got: %s\nWant: %s", got, tc.want)
		}
	}

}

func TestShowsQuery(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if sfx := "/shows/query?apikey=abc&limit=5&order=ASC"; !strings.HasSuffix(req.URL.String(), sfx) {
			t.Errorf("Expected request %s to have the suffix %s", req.URL.String(), sfx)
		}
		res.Write([]byte(`{"error_code":0,"error_message":null,"response":{"count":5,"data":[{"showid":1326251770,"showdate":"1982-12-07","artistid":2,"billed_as":"Space Antelope","link":"http:\/\/phish.net\/setlists\/trey-anastasio-december-07-1982-the-taft-school-watertown-ct-usa.html","location":"Watertown, CT, USA","venue":"The Taft School","setlistnotes":"This list is likely incomplete, and the date may be incorrect. This is the only show by Space Antelope documented by a recording in circulation, though there were likely other gigs. The Jam that preceded Walk on the Wild Side&nbsp;contained part of what would become the <a href=\"http:\/\/phish.net\/song\/arrival\/history\">Arrival<\/a>&nbsp;segment of <a href=\"http:\/\/phish.net\/song\/fluffs-travels\/history\">Fluff&#39;s Travels<\/a>.","venueid":1140,"tourid":61,"tourname":"Not Part of a Tour","tour_when":"No Tour","artistlink":"http:\/\/phish.net\/setlists\/trey"},{"showid":1251168326,"showdate":"1983-10-30","artistid":1,"billed_as":"Phish","link":"http:\/\/phish.net\/setlists\/phish-october-30-1983-harris-millis-cafeteria-university-of-vermont-burlington-vt-usa.html","location":"Burlington, VT, USA","venue":"Harris-Millis Cafeteria, UVM","setlistnotes":"Throughout most of Phish history this was understood to have been the date of the first Phish show. The band believed this to be true as late as 1998, when on October 30 they celebrated their &ldquo;15th anniversary.&rdquo; Later research, however, revealed this to be incorrect, and that the correct date of this first show &ndash; commonly referred to as the &ldquo;Thriller&rdquo; show or a &ldquo;Halloween Dance&rdquo; &ndash; is December 2, 1983.","venueid":7,"tourid":61,"tourname":"Not Part of a Tour","tour_when":"No Tour","artistlink":"http:\/\/phish.net\/setlists\/phish"},{"showid":1251253100,"showdate":"1983-12-02","artistid":1,"billed_as":"Phish","link":"http:\/\/phish.net\/setlists\/phish-december-02-1983-harris-millis-cafeteria-university-of-vermont-burlington-vt-usa.html","location":"Burlington, VT, USA","venue":"Harris-Millis Cafeteria, UVM","setlistnotes":"Trey, Mike, Fish, and Jeff Holdsworth recall being billed as &ldquo;Blackwood Convention&rdquo; for this show (though no one is certain what band name was used), which is believed to be their first public gig together. The band was short on equipment, so a hockey stick was used as a microphone stand. Between sets, the DJ spun some Michael Jackson and Trey drummed along to the album. The house music (which included more Michael Jackson) was presumably turned up after Fire to drown out the band. The setlist may be incomplete, though, as the master recording contains nothing after Trey&rsquo;s sarcastic comments about Michael Jackson following Fire. All songs were, of course, Phish debuts. Back in Black was teased before Scarlet Begonias. While this show is often billed as an ROTC Halloween Dance that took place on October 30, 1983, this is incorrect. The master copy of the recording of this show, as unearthed by Phish archivist Kevin Shapiro, contains a handwritten note that pegs the date as December 2, 1983. In discussions with Kevin, band members confirmed that they recall rehearsing for this show over the Thanksgiving Break, and that the show was a Christmas semi-formal. Also, it was not an ROTC-sponsored event; it was a dorm dance in a predominantly ROTC dorm (Mike&rsquo;s dorm at the time).","venueid":7,"tourid":1,"tourname":"1983 Tour","tour_when":"1983","artistlink":"http:\/\/phish.net\/setlists\/phish"},{"showid":1251253531,"showdate":"1983-12-03","artistid":1,"billed_as":"Phish","link":"http:\/\/phish.net\/setlists\/phish-december-03-1983-marsh-austin-tupper-dormitory-university-of-vermont-burlington-vt-usa.html","location":"Burlington, VT, USA","venue":"Marsh \/ Austin \/ Tupper Dormitory, University of Vermont","setlistnotes":"This show, played by Trey, Mike, Fish, and Jeff, may have been billed as &ldquo;Blackwood Convention.&rdquo; This date is believed to be correct but, due to a lack of records, the exact date cannot be ascertained. &nbsp; &nbsp;","venueid":272,"tourid":1,"tourname":"1983 Tour","tour_when":"1983","artistlink":"http:\/\/phish.net\/setlists\/phish"},{"showid":1331863583,"showdate":"1984-04-28","artistid":-1,"billed_as":"Dangerous Grapes","link":"http:\/\/phish.net\/setlists\/guest-appearance-april-28-1984-slade-hall-university-of-vermont-burlington-vt-usa.html","location":"Burlington, VT, USA","venue":"Slade Hall, University of Vermont","setlistnotes":"No setlist is known for this instance of Mike and Fish playing with the Dangerous Grapes.","venueid":246,"tourid":61,"tourname":"Not Part of a Tour","tour_when":"No Tour","artistlink":"http:\/\/phish.net\/setlists\/guest"}]}}`))
	}))
	defer func() { testServer.Close() }()

	c := NewClient("abc", WithBaseUrl(testServer.URL))
	resp, err := c.ShowsQuery(&ShowsQueryRequest{Order: "ASC", Limit: 5})
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}
	if resp.Response.Count != 5 {
		t.Errorf("Expected a count of 5, but got: %d", resp.Response.Count)
	}
	if len(resp.Response.Data) != 5 {
		t.Errorf("Expected 5 shows in data, but got %v", resp.Response.Data)
	}
}

func TestSetlistsGet(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if sfx := "/setlists/get?apikey=abc&showid=1250878246"; !strings.HasSuffix(req.URL.String(), sfx) {
			t.Errorf("Expected request %s to have the suffix %s", req.URL.String(), sfx)
		}
		res.Write([]byte(`{"error_code":0,"error_message":null,"response":{"count":1,"data":[{"showid":1250878246,"showdate":"1992-04-23","short_date":"04\/23\/1992","long_date":"Thursday 04\/23\/1992","relative_date":"27 years ago","url":"http:\/\/phish.net\/setlists\/phish-april-23-1992-oz-nightclub-seattle-wa-usa.html","gapchart":"http:\/\/phish.net\/setlists\/gap-chart\/phish-april-23-1992-oz-nightclub-seattle-wa-usa.html","artist":"<a href='http:\/\/phish.net\/setlists\/phish'>Phish<\/a>","artistid":1,"venueid":142,"venue":"<a href=\"http:\/\/phish.net\/venue\/142\/Oz_Nightclub\">Oz Nightclub<\/a>","location":"Seattle, WA, USA","setlistdata":"<p><span class='set-label'>Set 1<\/span>: <a href='http:\/\/phish.net\/song\/cavern' class='setlist-song'>Cavern<\/a> > <a href='http:\/\/phish.net\/song\/the-curtain' class='setlist-song'>The Curtain<\/a> > <a href='http:\/\/phish.net\/song\/split-open-and-melt' class='setlist-song'>Split Open and Melt<\/a>, <a href='http:\/\/phish.net\/song\/uncle-pen' class='setlist-song'>Uncle Pen<\/a>, <a href='http:\/\/phish.net\/song\/guelah-papyrus' class='setlist-song'>Guelah Papyrus<\/a>, <a href='http:\/\/phish.net\/song\/the-squirming-coil' class='setlist-song'>The Squirming Coil<\/a> > <a href='http:\/\/phish.net\/song\/llama' class='setlist-song'>Llama<\/a>, <a href='http:\/\/phish.net\/song\/bouncing-around-the-room' class='setlist-song'>Bouncing Around the Room<\/a>, <a href='http:\/\/phish.net\/song\/its-ice' class='setlist-song'>It's Ice<\/a>, <a href='http:\/\/phish.net\/song\/i-didnt-know' class='setlist-song'>I Didn't Know<\/a>, <a href='http:\/\/phish.net\/song\/possum' class='setlist-song'>Possum<\/a><sup title=\"Simpsons, Aw Fuck!, Oom Pa Pa, and All Fall Down signals in intro. Simpsons signal at the end.\">[1]<\/sup><\/p><p><span class='set-label'>Set 2<\/span>: <a href='http:\/\/phish.net\/song\/the-landlady' class='setlist-song'>The Landlady<\/a>, <a href='http:\/\/phish.net\/song\/poor-heart' class='setlist-song'>Poor Heart<\/a> > <a href='http:\/\/phish.net\/song\/mikes-song' class='setlist-song'>Mike's Song<\/a> > <a href='http:\/\/phish.net\/song\/i-am-hydrogen' class='setlist-song'>I Am Hydrogen<\/a> > <a href='http:\/\/phish.net\/song\/weekapaug-groove' class='setlist-song'>Weekapaug Groove<\/a>, <a href='http:\/\/phish.net\/song\/the-lizards' class='setlist-song'>The Lizards<\/a>, <a href='http:\/\/phish.net\/song\/nicu' class='setlist-song'>NICU<\/a>, <a href='http:\/\/phish.net\/song\/horn' class='setlist-song'>Horn<\/a> > <a href='http:\/\/phish.net\/song\/tweezer' class='setlist-song'>Tweezer<\/a> > <a href='http:\/\/phish.net\/song\/fee' class='setlist-song'>Fee<\/a><sup title=\"Lyrics changed to \"have a cup of espresso.\"\">[2]<\/sup> -> <a href='http:\/\/phish.net\/song\/maze' class='setlist-song'>Maze<\/a>, <a href='http:\/\/phish.net\/song\/cold-as-ice' class='setlist-song'>Cold as Ice<\/a> > <a href='http:\/\/phish.net\/song\/cracklin-rosie' class='setlist-song'>Cracklin' Rosie<\/a> > <a href='http:\/\/phish.net\/song\/cold-as-ice' class='setlist-song'>Cold as Ice<\/a>, <a href='http:\/\/phish.net\/song\/golgi-apparatus' class='setlist-song'>Golgi Apparatus<\/a><\/p><p><span class='set-label'>Encore<\/span>: <a href='http:\/\/phish.net\/song\/sleeping-monkey' class='setlist-song'>Sleeping Monkey<\/a> > <a href='http:\/\/phish.net\/song\/tweezer-reprise' class='setlist-song'>Tweezer Reprise<\/a><p class='setlist-footer'>[1] Simpsons, Aw Fuck!, Oom Pa Pa, and All Fall Down signals in intro. Simpsons signal at the end.<br>[2] Lyrics changed to \"have a cup of espresso.\"<br><\/p>","setlistnotes":"Fish&nbsp;was in the process of introducing Rosie as being by one of his favorite composers when feedback caused him to say &quot;Jesus Christ,&quot; leading to some amusing banter. The Possum intro contained a full band Yield Not to Temptation tease and Simpsons, Aw Fuck!, Oom Pa Pa, and All Fall Down signals. Possum also contained a Dixie tease as well as a Simpsons signal at the end. Landlady was preceded by a Random Laugh signal. Lizards was played for &quot;Liz.&quot; Fee&#39;s lyrics were changed to &quot;have a cup of espresso.&quot;<br>via <a href=\"http:\/\/phish.net\">phish.net<\/a>","rating":"2.8750"}]}}`))
	}))
	defer func() { testServer.Close() }()

	c := NewClient("abc", WithBaseUrl(testServer.URL))
	resp, err := c.SetlistsGet(&SetlistsGetRequest{ShowId: 1250878246})
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}
	if resp.Response.Count != 1 {
		t.Errorf("Expected a count of 1, but got: %d", resp.Response.Count)
	}
	if len(resp.Response.Data) != 1 {
		t.Errorf("Expected 1 setlist in data, but got %v", resp.Response.Data)
	}
}
