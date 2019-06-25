package badge

import (
	"net/http"
	"strings"

	svg "github.com/ajstarks/svgo"
	"github.com/argoproj/argo-cd/pkg/client/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

  "github.com/argoproj/argo-cd/pkg/client/clientset/versioned"
)

//NewHandler creates handler serving to do api/badge endpoint
func NewHandler(appClientset versioned.Interface, namespace string) http.Handler {
	return &Handler{appClientset: appClientset, namespace: namespace}
}

//Handler used to get application in order to access health/sync
type Handler struct {
	namespace    string
	appClientset versioned.Interface
}

//ServeHTTP returns badge with health and sync status for application
//(or an error badge if wrong query or application name is given)
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//Sample url: http://localhost:8080/api/badge?name=123
	keys, ok := r.URL.Query()["name"]

	//error when you add another query after name = and do like http://localhost:8080/api/badge?name=123/hadfkjajdhj

	pageWidth := 2000
	pageHeight := 2000
	xStart := 0
	yStart := 0
	badgeHeight := 25
	//badgeCurve is the rx/ry value for the round rectangle for each badge
	badgeCurve := 2
	//xHealthText is x pos where text for health badge and edge case badges start
	xHealthText := 3

	textFormat := "font-size:11;fill:black"
	yText := 17

	if !ok || len(keys[0]) < 1 {
		svgOne := svg.New(w)
		svgOne.Start(pageWidth, pageHeight)
		svgOne.Roundrect(xStart, yStart, 164, badgeHeight, badgeCurve, badgeCurve, "fill:rgb(200,321,233);stroke:black")
		svgOne.Text(xHealthText, yText, "Param 'name=' missing or left blank", textFormat)
		svgOne.End()
		return
	}

	key := keys[0]
	//if another query is added after the appplication name and is separated by a / this will make sure it only looks at
	//what is between the name= and / and will open the applicaion by that name
	q := strings.Split(key, "/")
	key = q[0]
	app, err := h.appClientset.ArgoprojV1alpha1().Applications(h.namespace).Get(key, v1.GetOptions{})
	if err != nil {
		notFoundBadgeLength := len("Application"+key+" not found") * 6
		svgTwo := svg.New(w)
		svgTwo.Start(pageWidth, pageHeight)
		svgTwo.Roundrect(xStart, yStart, notFoundBadgeLength, badgeHeight, badgeCurve, badgeCurve, "fill:rgb(252, 197, 0);stroke:black")
		svgTwo.Text(xHealthText, yText, "Application '"+key+"' not found", textFormat)
		svgTwo.End()
		return
	}

	health := app.Status.Health.Status
	status := app.Status.Sync.Status
	//healthBadgeLength is where the x value sync badge starts along with being the length of the health badge
	healthBadgeLength := len(health) * 7
	syncBadgeLength := len(status) * 7
	syncTextStart := healthBadgeLength + 3
	svgThree := svg.New(w)
	svgThree.Start(pageWidth, pageHeight)

	switch health {
	case appv1.HealthStatusHealthy:
		svgThree.Roundrect(xStart, yStart, healthBadgeLength, badgeHeight, badgeCurve, badgeCurve, "fill:rgb(127,255,131);stroke:black")
		svgThree.Text(xHealthText, yText, health, textFormat)
	case appv1.HealthStatusProgressing:
		svgThree.Roundrect(xStart, yStart, healthBadgeLength, badgeHeight, badgeCurve, badgeCurve, "fill:rgb(255,251,92);stroke:black")
		svgThree.Text(xHealthText, yText, health, textFormat)
	case appv1.HealthStatusSuspended:
		svgThree.Roundrect(xStart, yStart, healthBadgeLength, badgeHeight, badgeCurve, badgeCurve, "fill:rgb(255,145,0);stroke:black")
		svgThree.Text(xHealthText, yText, health, textFormat)
	case appv1.HealthStatusDegraded:
		svgThree.Roundrect(xStart, yStart, healthBadgeLength, badgeHeight, badgeCurve, badgeCurve, "fill:rgb(109,202,205);stroke:black")
		svgThree.Text(xHealthText, yText, health, textFormat)
	case appv1.HealthStatusMissing:
		svgThree.Roundrect(xStart, yStart, healthBadgeLength, badgeHeight, badgeCurve, badgeCurve, "fill:rgb(255,36,36);stroke:black")
		svgThree.Text(xHealthText, yText, health, textFormat)
	default:
		svgThree.Roundrect(xStart, yStart, healthBadgeLength, badgeHeight, badgeCurve, badgeCurve, "fill:rgb(178,102,255);stroke:black")
		svgThree.Text(xHealthText, yText, health, textFormat)
	}
	switch status {
	case appv1.SyncStatusCodeSynced:
		svgThree.Roundrect(healthBadgeLength, yStart, syncBadgeLength, badgeHeight, badgeCurve, badgeCurve, "fill:rgb(0,204,0);stroke:black")
		svgThree.Text(syncTextStart, yText, string(appv1.SyncStatusCodeSynced), textFormat)
	case appv1.SyncStatusCodeOutOfSync:
		svgThree.Roundrect(healthBadgeLength, yStart, syncBadgeLength, badgeHeight, badgeCurve, badgeCurve, "fill:rgb(255,57,57);stroke:black")
		svgThree.Text(syncTextStart, yText, string(appv1.SyncStatusCodeOutOfSync), textFormat)
	default:
		svgThree.Roundrect(healthBadgeLength, yStart, syncBadgeLength, badgeHeight, badgeCurve, badgeCurve, "fill:rgb(209,155,177);stroke:black")
		svgThree.Text(syncTextStart, yText, string(appv1.SyncStatusCodeUnknown), textFormat)
	}
	svgThree.End()
}
