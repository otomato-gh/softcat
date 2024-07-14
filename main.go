package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	models "github.com/otomato/softcat/models"
)

// postComponents adds a new component to the list.
func postComponents(c *gin.Context) {
	var newComponent models.Component

	if err := c.BindJSON(&newComponent); err != nil {
		log.Println("error binding component", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := models.PostComponents(newComponent); err != nil {
		log.Println("error inserting component:", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.IndentedJSON(http.StatusCreated, newComponent)
}

// postTeams adds a new team to the list.
func postTeams(c *gin.Context) {
	var newTeam models.Team
	if err := c.BindJSON(&newTeam); err != nil {
		return
	}
	models.PostTeams(newTeam)
	c.IndentedJSON(http.StatusCreated, newTeam)
}

// getComponents responds with the list of all components as JSON.
func getComponents(c *gin.Context) {
	components := models.GetComponents()
	c.IndentedJSON(http.StatusOK, components)
}

// getTeams responds with the list of all teams as JSON.
func getTeams(c *gin.Context) {
	teams := models.GetTeams()
	c.IndentedJSON(http.StatusOK, teams)
}

// getComponentsByTeam responds with the list of all components for a given team as JSON.
func getComponentsByTeam(c *gin.Context) {
	teamID, _ := strconv.Atoi(c.Param("teamID"))
	teamComponents := models.GetComponentsByTeam(teamID)

	c.IndentedJSON(http.StatusOK, teamComponents)
}

// getComponentsById responds with the component with the given ID as JSON.
func getComponentsById(c *gin.Context) {
	componentID := c.Param("componentID")
	component, err := models.GetComponentByID(componentID)
	if err == nil {
		c.IndentedJSON(http.StatusOK, component)
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "component not found"})
}

// getTeamsById responds with the team with the given ID as JSON.
func getTeamById(c *gin.Context) {
	teamID, _ := strconv.Atoi(c.Param("teamID"))
	team, err := models.GetTeamByID(teamID)
	if err == nil {
		c.IndentedJSON(http.StatusOK, team)
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "team not found"})
}

func generateIdenticons() {
	for {
		var components []models.Component = models.GetComponents()
		for _, component := range components {
			img, _ := models.GetImgByComponent(component.ID)
			if img == nil {
				log.Println("generating identicon for component", component.ID)
				res, err := http.Post("http://localhost:8081",
					"application/json",
					bytes.NewBuffer([]byte(`{"id": `+strconv.Itoa(component.ID)+`, "name": "`+component.Name+`"}`)))
				if err != nil {
					log.Println("failed to generate identicon for component", component.ID, err)
				} else if res.StatusCode != http.StatusOK {
					log.Println("failed to generate identicon for component", component.ID)
				} else {
					identicon, _ := io.ReadAll(res.Body)
					img := models.Image{ID: component.ID, Image: identicon}
					models.PostImg(img)
				}
			}
		}
		log.Println(time.Now().UTC())
		time.Sleep(10000 * time.Millisecond)
	}
}

func main() {

	models.ConnectDB()

	//generate identicons in the background
	go generateIdenticons()
	// run gin server
	r := gin.Default()
	pprof.Register(r)
	r.GET("/components", getComponents)
	r.POST("/components", postComponents)
	r.GET("/components/:componentID", getComponentsById)
	r.GET("/teams", getTeams)
	r.POST("/teams", postTeams)
	r.GET("/teams/:teamID", getTeamById)
	r.GET("/teams/:teamID/components", getComponentsByTeam)
	r.Run(":8080")
}
