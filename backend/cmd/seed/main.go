package main

import (
	"backend/internal/database"
	"backend/internal/domain/domain"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	// Load .env file from backend root
	// TODO Call API to get university
	if err := godotenv.Load(); err != nil {
		log.Printf("⚠️ Warning: .env file not found, using system env vars")
	}

	db := database.ConnectDB()

	fmt.Println("🌱 Starting database seeding...")

	// seedInstitutes(db)
	// seedFieldOfStudies(db)
	seedSkills(db)

	fmt.Println("✅ Seeding completed!")
}

func seedSkills(db *gorm.DB) {
	skills := []string{
		// Languages
		"Go", "Golang", "Python", "JavaScript", "TypeScript", "Java", "C#", "C++", "C", "Rust", "Swift", "Kotlin", "PHP", "Ruby", "Dart", "Scala", "Elixir", "Haskell", "Lua", "Perl", "R", "Shell", "SQL", "HTML", "CSS", "Assembly", "Matlab", "Groovy", "Objective-C", "VBA", "Visual Basic",

		// Frontend Frameworks/Libraries
		"React", "Vue.js", "Angular", "Svelte", "Next.js", "Nuxt.js", "SolidJS", "Qwik", "jQuery", "Ember.js", "Backbone.js", "Preact", "Alpine.js", "Lit", "Stencil", "Tailwind CSS", "Bootstrap", "Material UI", "Chakra UI", "Ant Design", "Bulma", "Sass", "Less", "Styled Components",

		// Backend Frameworks/Libraries
		"Node.js", "Express.js", "NestJS", "Fastify", "Django", "Flask", "FastAPI", "Spring Boot", "Laravel", "Symfony", "Ruby on Rails", "ASP.NET Core", "Gin", "Fiber", "Echo", "Chi", "Beego", "Revel", "Phoenix", "Ktor", "Vapor",

		// Mobile
		"Flutter", "React Native", "SwiftUI", "Jetpack Compose", "Xamarin", "Ionic", "Cordova", "Capacitor", "Expo", "Unity", "Unreal Engine", "Godot",

		// Database
		"PostgreSQL", "MySQL", "MariaDB", "SQLite", "MongoDB", "Redis", "Cassandra", "Elasticsearch", "DynamoDB", "Firestore", "CouchDB", "Neo4j", "Oracle", "Microsoft SQL Server", "CockroachDB", "TiDB", "ClickHouse", "InfluxDB", "Prometheus", "Supabase", "Firebase",

		// DevOps & Cloud
		"Docker", "Kubernetes", "AWS", "Google Cloud Platform", "Azure", "Terraform", "Ansible", "Jenkins", "GitLab CI", "GitHub Actions", "CircleCI", "Travis CI", "Nginx", "Apache", "Linux", "Ubuntu", "CentOS", "Debian", "Bash", "PowerShell", "Vagrant", "OpenShift", "Heroku", "Netlify", "Vercel", "DigitalOcean", "Linode", "Cloudflare",

		// Tools & Others
		"Git", "GitHub", "GitLab", "Bitbucket", "Jira", "Confluence", "Trello", "Notion", "Slack", "Discord", "Zoom", "Microsoft Teams", "Postman", "Insomnia", "Swagger", "OpenAPI", "GraphQL", "gRPC", "WebSocket", "WebRTC", "Socket.io", "Kafka", "RabbitMQ", "ActiveMQ", "ZeroMQ", "NATS", "Redis Pub/Sub", "Celery", "Bull, Sidekiq",

		// AI/ML
		"TensorFlow", "PyTorch", "Keras", "Scikit-learn", "Pandas", "NumPy", "Matplotlib", "Seaborn", "OpenCV", "NLTK", "Spacy", "Hugging Face", "OpenAI API", "LangChain",
	}

	count := 0
	for _, name := range skills {
		result := db.FirstOrCreate(&domain.Skill{}, domain.Skill{Name: name})
		if result.Error != nil {
			log.Printf("❌ Error seeding skill %s: %v", name, result.Error)
		} else if result.RowsAffected > 0 {
			count++
		}
	}
	log.Printf("✅ Seeded %d skills (total list: %d)", count, len(skills))
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
