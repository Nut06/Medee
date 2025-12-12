package main

import (
	"backend/internal/database"
	"backend/internal/domain/domain"
	"fmt"
	"log"

	"gorm.io/gorm"
)

func main() {
	db := database.ConnectDB()

	fmt.Println("🌱 Starting database seeding...")

	seedInstitutes(db)
	seedFieldOfStudies(db)

	fmt.Println("✅ Seeding completed!")
}

func seedInstitutes(db *gorm.DB) {
	institutes := []domain.Institute{
		// Thai Universities
		{Name: "King Mongkut's Institute of Technology Ladkrabang"},
		{Name: "Chulalongkorn University"},
		{Name: "Mahidol University"},
		{Name: "Thammasat University"},
		{Name: "Kasetsart University"},
		{Name: "Chiang Mai University"},
		{Name: "Khon Kaen University"},
		{Name: "Prince of Songkla University"},
		{Name: "Rajamangala University of Technology"},
		{Name: "Silpakorn University"},

		// International Universities
		{Name: "Stanford University"},
		{Name: "Massachusetts Institute of Technology"},
		{Name: "Harvard University"},
		{Name: "University of California, Berkeley"},
		{Name: "University of Oxford"},
		{Name: "University of Cambridge"},
		{Name: "National University of Singapore"},
		{Name: "Nanyang Technological University"},
		{Name: "University of Tokyo"},
		{Name: "Seoul National University"},
		{Name: "Peking University"},
		{Name: "Tsinghua University"},
		{Name: "ETH Zurich"},
		{Name: "Imperial College London"},
		{Name: "University of Toronto"},
		{Name: "University of Melbourne"},
		{Name: "Carnegie Mellon University"},
		{Name: "Cornell University"},
		{Name: "Yale University"},
		{Name: "Princeton University"},
	}

	count := 0
	for _, institute := range institutes {
		result := db.FirstOrCreate(&institute, domain.Institute{Name: institute.Name})
		if result.Error != nil {
			log.Printf("❌ Error seeding institute %s: %v", institute.Name, result.Error)
		} else if result.RowsAffected > 0 {
			count++
		}
	}

	log.Printf("✅ Seeded %d institutes (total: %d)", count, len(institutes))
}

func seedFieldOfStudies(db *gorm.DB) {
	fields := []domain.FieldOfStudy{
		// Computer Science & IT
		{Name: "Computer Science"},
		{Name: "Computer Engineering"},
		{Name: "Software Engineering"},
		{Name: "Information Technology"},
		{Name: "Data Science"},
		{Name: "Artificial Intelligence"},
		{Name: "Machine Learning"},
		{Name: "Cybersecurity"},
		{Name: "Information Systems"},
		{Name: "Computer Graphics"},

		// Engineering
		{Name: "Electrical Engineering"},
		{Name: "Mechanical Engineering"},
		{Name: "Civil Engineering"},
		{Name: "Chemical Engineering"},
		{Name: "Industrial Engineering"},
		{Name: "Aerospace Engineering"},
		{Name: "Biomedical Engineering"},
		{Name: "Environmental Engineering"},
		{Name: "Materials Engineering"},
		{Name: "Robotics Engineering"},

		// Business
		{Name: "Business Administration"},
		{Name: "Economics"},
		{Name: "Accounting"},
		{Name: "Finance"},
		{Name: "Marketing"},
		{Name: "Management"},
		{Name: "International Business"},
		{Name: "Entrepreneurship"},
		{Name: "Human Resource Management"},

		// Science
		{Name: "Physics"},
		{Name: "Chemistry"},
		{Name: "Biology"},
		{Name: "Mathematics"},
		{Name: "Statistics"},
		{Name: "Biotechnology"},
		{Name: "Environmental Science"},

		// Arts & Social Sciences
		{Name: "Psychology"},
		{Name: "Sociology"},
		{Name: "Political Science"},
		{Name: "Communication Arts"},
		{Name: "Graphic Design"},
		{Name: "Architecture"},
		{Name: "Interior Design"},

		// Health Sciences
		{Name: "Medicine"},
		{Name: "Nursing"},
		{Name: "Pharmacy"},
		{Name: "Public Health"},
		{Name: "Dentistry"},

		// Other
		{Name: "Law"},
		{Name: "Education"},
		{Name: "Journalism"},
		{Name: "Digital Media"},
		{Name: "Game Development"},
	}

	count := 0
	for _, field := range fields {
		result := db.FirstOrCreate(&field, domain.FieldOfStudy{Name: field.Name})
		if result.Error != nil {
			log.Printf("❌ Error seeding field %s: %v", field.Name, result.Error)
		} else if result.RowsAffected > 0 {
			count++
		}
	}

	log.Printf("✅ Seeded %d field of studies (total: %d)", count, len(fields))
}
