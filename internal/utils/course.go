package utils

import (
	"chaoxing/internal/models"
	"strings"
)

func ParseCourse(data string) []models.CourseType {
	var courses []models.CourseType
	i := 0 // 全局索引

	for {
		courseIndex := strings.Index(data[i:], "course_")
		if courseIndex == -1 {
			break
		}

		i += courseIndex
		endOfCourseId := strings.Index(data[i+len("course_"):], "_")
		if endOfCourseId == -1 {
			break
		}

		endOfCourseId += i + len("course_")
		courseId := data[i+len("course_") : endOfCourseId]

		classIdStart := endOfCourseId + 1
		classIdEnd := strings.Index(data[classIdStart:], `"`)
		if classIdEnd == -1 {
			break
		}

		classIdEnd += classIdStart
		classId := data[classIdStart:classIdEnd]

		courses = append(courses, models.CourseType{
			CourseID: courseId,
			ClassID:  classId,
		})

		i = classIdEnd + 1
	}

	return courses
}
