package web

import (
	structs "github.com/whoismissing/gizmo/structs"

	"strconv"
)

func GenerateScoreboardHTML(teams []structs.Team) string {
	html := `<!DOCTYPE html> 
<html> 
<head>
<meta http-equiv="refresh" content="5" />
</head>
<body> 
<table style="width:100%"> 
<tr> 
<th>Team</th>
`

	services := teams[0].Services
	for i := 0; i < len(services); i++ {
		html += "<th>" + services[i].Name + "</th>\n"
	}

	for i := 0; i < len(teams); i++ {
		html += "<tr>\n"
		html += "<td>Team " + strconv.Itoa(int(teams[i].TeamID)) + "</td>\n"
		services := teams[i].Services
		for j := 0; j < len(services); j++ {
			html += "<td style=\"background-color:"
			status := services[j].Status
			if status == true {
				html += "green"
			} else {
				html += "red"
			}
			html += "\">"
            html += services[j].HostIP
            html += "&emsp;" // 4 spaces in html

            top := len(services[j].PrevStatuses) - 1
            /* if service.PrevStatuses is empty */
            if top < 0 {

            } else {
                html += services[j].PrevStatuses[top].Time.Format("2006-01-02 15:04:05")
            }

            html += "</td>\n"
		}
		html += "</tr>\n"
	}

	html += "</table>\n"
	html += "</body>\n"
	html += "</html>\n"

	return html
}
