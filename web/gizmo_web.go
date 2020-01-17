package web

import (
	structs "github.com/whoismissing/gizmo/structs"

	"strconv"
)

func GenerateScoreboardHTML(teams []structs.Team) string {
	html := `<!DOCTYPE html> 
<html> 
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
			html += "\"></td>\n"
		}
		html += "</tr>\n"
	}

	html += "</table>\n"
	html += "</body>\n"
	html += "</html>\n"

	return html
}
