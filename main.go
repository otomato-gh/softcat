package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// album represents data about a record album.
type component struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Team     int    `json:"team"`
	Language string `json:"language"`
}

type team struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var components = []component{
	{ID: "1", Name: "Component 1", Type: "Service", Team: 1, Language: "Java"},
	{ID: "2", Name: "Component 2", Type: "Library", Team: 1, Language: "Golang"},
	{ID: "3", Name: "Component 3", Type: "Service", Team: 2, Language: "Python"},
	{ID: "4", Name: "Component 4", Type: "Library", Team: 2, Language: "Golang"},
	{ID: "5", Name: "Component 5", Type: "Data Pipeline", Team: 3, Language: "Java"},
}

var teams = []team{
	{ID: 1, Name: "Team 1"},
	{ID: 2, Name: "Team 2"},
	{ID: 3, Name: "Team 3"},
}

// postComponents adds a new component to the list.
func postComponents(c *gin.Context) {
	var newComponent component
	if err := c.BindJSON(&newComponent); err != nil {
		return
	}
	components = append(components, newComponent)
	c.IndentedJSON(http.StatusCreated, newComponent)
}

// postTeams adds a new team to the list.
func postTeams(c *gin.Context) {
	var newTeam team
	if err := c.BindJSON(&newTeam); err != nil {
		return
	}
	teams = append(teams, newTeam)
	c.IndentedJSON(http.StatusCreated, newTeam)
}

// getComponents responds with the list of all components as JSON.
func getComponents(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, components)
}

// getTeams responds with the list of all teams as JSON.
func getTeams(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, teams)
}

// getComponentsByTeam responds with the list of all components for a given team as JSON.
func getComponentsByTeam(c *gin.Context) {
	teamID := c.Param("teamID")
	var teamComponents []component
	for _, component := range components {
		if strconv.Itoa(component.Team) == teamID {
			teamComponents = append(teamComponents, component)
		}
	}
	c.IndentedJSON(http.StatusOK, teamComponents)
}

// getComponentsById responds with the component with the given ID as JSON.
func getComponentsById(c *gin.Context) {
	componentID := c.Param("componentID")
	for _, component := range components {
		if component.ID == componentID {
			c.IndentedJSON(http.StatusOK, component)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "component not found"})
}

// getTeamsById responds with the team with the given ID as JSON.
func getTeamsById(c *gin.Context) {
	teamID := c.Param("teamID")
	for _, team := range teams {
		if strconv.Itoa(team.ID) == teamID {
			c.IndentedJSON(http.StatusOK, team)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "team not found"})
}

func main() {
	dbconn := new(DBconn)
	var err error
	dbconn.DB, err = ConnectDB()
	if err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}

	r := gin.Default()
	r.GET("/components", getComponents)
	r.POST("/components", postComponents)
	r.GET("/components/:componentID", getComponentsById)
	r.GET("/teams", getTeams)
	r.POST("/teams", postTeams)
	r.GET("/teams/:teamID", getTeamsById)
	r.GET("/teams/:teamID/components", getComponentsByTeam)
	r.Run(":8080")
}
