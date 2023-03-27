package provider

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sivaramsajeev/terraform-provider-student/api/client"
	"github.com/sivaramsajeev/terraform-provider-student/api/server"
)

func validateName(v interface{}, k string) (ws []string, es []error) {
	var errs []error
	var warns []string
	value, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("Expected name to be string"))
		return warns, errs
	}
	whiteSpace := regexp.MustCompile(`\s+`)
	if whiteSpace.Match([]byte(value)) {
		errs = append(errs, fmt.Errorf("name cannot contain whitespace. Got %s", value))
		return warns, errs
	}
	return warns, errs
}

func resourceStudent() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the student & it's also the unique ID",
				ForceNew:     true,
				ValidateFunc: validateName,
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A description of the student",
			},
			"subjects": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "An optional list of subjects, represented as a key, value pair",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		Create: resourceCreateStudent,
		Read:   resourceReadStudent,
		Update: resourceUpdateStudent,
		Delete: resourceDeleteStudent,
		Exists: resourceExistsStudent,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateStudent(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	tfSubs := d.Get("subjects").(*schema.Set).List()
	subs := make([]string, len(tfSubs))
	for i, sub := range tfSubs {
		subs[i] = sub.(string)
	}

	student := server.Student{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Subjects:    subs,
	}

	err := apiClient.NewStudent(&student)

	if err != nil {
		return err
	}
	d.SetId(student.Name) // Marks the attribute as the unique id for the resource
	return nil
}

func resourceReadStudent(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	studentId := d.Id()
	student, err := apiClient.GetStudent(studentId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
		} else {
			return fmt.Errorf("error finding Student with ID %s", studentId)
		}
	}

	d.SetId(student.Name) // Marks the attribute as the unique id for the resource
	d.Set("name", student.Name)
	d.Set("description", student.Description)
	d.Set("subjects", student.Subjects)
	return nil
}

func resourceUpdateStudent(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	tfSubs := d.Get("subjects").(*schema.Set).List()
	subs := make([]string, len(tfSubs))
	for i, sub := range tfSubs {
		subs[i] = sub.(string)
	}

	student := server.Student{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Subjects:    subs,
	}

	err := apiClient.UpdateStudent(&student)
	if err != nil {
		return err
	}
	return nil
}

func resourceDeleteStudent(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	studentId := d.Id()

	err := apiClient.DeleteStudent(studentId)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceExistsStudent(d *schema.ResourceData, m interface{}) (bool, error) {
	apiClient := m.(*client.Client)

	studentId := d.Id()
	_, err := apiClient.GetStudent(studentId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
