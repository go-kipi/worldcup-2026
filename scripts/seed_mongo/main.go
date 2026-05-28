package main

import (
	"context"
	"log"
	"time"

	"github.com/go-kipi/worldcup-2026/internal/config"
	"github.com/go-kipi/worldcup-2026/internal/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if cfg.MongoURI == "" {
		log.Fatal("MONGO_URI is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(cfg.MongoURI).
		SetBSONOptions(&options.BSONOptions{
			ObjectIDAsHexString: true,
		}))
	if err != nil {
		log.Fatalf("failed to connect to mongodb: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database(cfg.MongoDBName)

	seedData(ctx, db)
}

func seedData(ctx context.Context, db *mongo.Database) {
	// 1. AppSettings
	appSettingsCol := db.Collection("app_settings")
	appSettingsCol.Drop(ctx)
	_, err := appSettingsCol.InsertOne(ctx, models.AppSetting{
		ID:     "1",
		LockAt: time.Date(2026, 6, 11, 16, 0, 0, 0, time.UTC),
	})
	if err != nil {
		log.Printf("Failed to seed app settings: %v", err)
	}

	// 2. Teams
	teamsCol := db.Collection("teams")
	teamsCol.Drop(ctx)
	teams := []interface{}{
		models.Team{ID: "MEX", Name: "Mexico", GroupCode: "A", FlagEmoji: "🇲🇽"},
		models.Team{ID: "RSA", Name: "South Africa", GroupCode: "A", FlagEmoji: "🇿🇦"},
		models.Team{ID: "KOR", Name: "South Korea", GroupCode: "A", FlagEmoji: "🇰🇷"},
		models.Team{ID: "CZE", Name: "Czechia", GroupCode: "A", FlagEmoji: "🇨🇿"},
		models.Team{ID: "CAN", Name: "Canada", GroupCode: "B", FlagEmoji: "🇨🇦"},
		models.Team{ID: "BIH", Name: "Bosnia and Herzegovina", GroupCode: "B", FlagEmoji: "🇧🇦"},
		models.Team{ID: "QAT", Name: "Qatar", GroupCode: "B", FlagEmoji: "🇶🇦"},
		models.Team{ID: "SUI", Name: "Switzerland", GroupCode: "B", FlagEmoji: "🇨🇭"},
		models.Team{ID: "BRA", Name: "Brazil", GroupCode: "C", FlagEmoji: "🇧🇷"},
		models.Team{ID: "MAR", Name: "Morocco", GroupCode: "C", FlagEmoji: "🇲🇦"},
		models.Team{ID: "HAI", Name: "Haiti", GroupCode: "C", FlagEmoji: "🇭🇹"},
		models.Team{ID: "SCO", Name: "Scotland", GroupCode: "C", FlagEmoji: "🏴"},
		models.Team{ID: "USA", Name: "United States", GroupCode: "D", FlagEmoji: "🇺🇸"},
		models.Team{ID: "PAR", Name: "Paraguay", GroupCode: "D", FlagEmoji: "🇵🇾"},
		models.Team{ID: "AUS", Name: "Australia", GroupCode: "D", FlagEmoji: "🇦🇺"},
		models.Team{ID: "TUR", Name: "Türkiye", GroupCode: "D", FlagEmoji: "🇹🇷"},
		models.Team{ID: "GER", Name: "Germany", GroupCode: "E", FlagEmoji: "🇩🇪"},
		models.Team{ID: "CUW", Name: "Curaçao", GroupCode: "E", FlagEmoji: "🇨🇼"},
		models.Team{ID: "CIV", Name: "Côte d'Ivoire", GroupCode: "E", FlagEmoji: "🇨🇮"},
		models.Team{ID: "ECU", Name: "Ecuador", GroupCode: "E", FlagEmoji: "🇪🇨"},
		models.Team{ID: "NED", Name: "Netherlands", GroupCode: "F", FlagEmoji: "🇳🇱"},
		models.Team{ID: "JPN", Name: "Japan", GroupCode: "F", FlagEmoji: "🇯🇵"},
		models.Team{ID: "SWE", Name: "Sweden", GroupCode: "F", FlagEmoji: "🇸🇪"},
		models.Team{ID: "TUN", Name: "Tunisia", GroupCode: "F", FlagEmoji: "🇹🇳"},
		models.Team{ID: "BEL", Name: "Belgium", GroupCode: "G", FlagEmoji: "🇧🇪"},
		models.Team{ID: "EGY", Name: "Egypt", GroupCode: "G", FlagEmoji: "🇪🇬"},
		models.Team{ID: "IRN", Name: "Iran", GroupCode: "G", FlagEmoji: "🇮🇷"},
		models.Team{ID: "NZL", Name: "New Zealand", GroupCode: "G", FlagEmoji: "🇳🇿"},
		models.Team{ID: "ESP", Name: "Spain", GroupCode: "H", FlagEmoji: "🇪🇸"},
		models.Team{ID: "CPV", Name: "Cape Verde", GroupCode: "H", FlagEmoji: "🇨🇻"},
		models.Team{ID: "KSA", Name: "Saudi Arabia", GroupCode: "H", FlagEmoji: "🇸🇦"},
		models.Team{ID: "URU", Name: "Uruguay", GroupCode: "H", FlagEmoji: "🇺🇾"},
		models.Team{ID: "FRA", Name: "France", GroupCode: "I", FlagEmoji: "🇫🇷"},
		models.Team{ID: "SEN", Name: "Senegal", GroupCode: "I", FlagEmoji: "🇸🇳"},
		models.Team{ID: "IRQ", Name: "Iraq", GroupCode: "I", FlagEmoji: "🇮🇶"},
		models.Team{ID: "NOR", Name: "Norway", GroupCode: "I", FlagEmoji: "🇳🇴"},
		models.Team{ID: "ARG", Name: "Argentina", GroupCode: "J", FlagEmoji: "🇦🇷"},
		models.Team{ID: "ALG", Name: "Algeria", GroupCode: "J", FlagEmoji: "🇩🇿"},
		models.Team{ID: "AUT", Name: "Austria", GroupCode: "J", FlagEmoji: "🇦🇹"},
		models.Team{ID: "JOR", Name: "Jordan", GroupCode: "J", FlagEmoji: "🇯🇴"},
		models.Team{ID: "POR", Name: "Portugal", GroupCode: "K", FlagEmoji: "🇵🇹"},
		models.Team{ID: "COD", Name: "DR Congo", GroupCode: "K", FlagEmoji: "🇨🇩"},
		models.Team{ID: "UZB", Name: "Uzbekistan", GroupCode: "K", FlagEmoji: "🇺🇿"},
		models.Team{ID: "COL", Name: "Colombia", GroupCode: "K", FlagEmoji: "🇨🇴"},
		models.Team{ID: "ENG", Name: "England", GroupCode: "L", FlagEmoji: "🏴"},
		models.Team{ID: "CRO", Name: "Croatia", GroupCode: "L", FlagEmoji: "🇭🇷"},
		models.Team{ID: "GHA", Name: "Ghana", GroupCode: "L", FlagEmoji: "🇬🇭"},
		models.Team{ID: "PAN", Name: "Panama", GroupCode: "L", FlagEmoji: "🇵🇦"},
	}
	_, err = teamsCol.InsertMany(ctx, teams)
	if err != nil {
		log.Printf("Failed to seed teams: %v", err)
	}

	// 3. Knockout Slots
	slotsCol := db.Collection("knockout_slots")
	slotsCol.Drop(ctx)
	slots := []interface{}{
		models.KnockoutSlot{Round: "R32", SlotCode: "R32-01", Label: ptr("Round of 32 - Match 1"), HomeLabel: ptr("1A"), AwayLabel: ptr("3C/D/E"), Kickoff: parseDate("2026-06-28 16:00:00")},
		models.KnockoutSlot{Round: "R32", SlotCode: "R32-02", Label: ptr("Round of 32 - Match 2"), HomeLabel: ptr("1B"), AwayLabel: ptr("3A/C/D"), Kickoff: parseDate("2026-06-28 20:00:00")},
		models.KnockoutSlot{Round: "R32", SlotCode: "R32-03", Label: ptr("Round of 32 - Match 3"), HomeLabel: ptr("1C"), AwayLabel: ptr("3B/E/F"), Kickoff: parseDate("2026-06-29 16:00:00")},
		models.KnockoutSlot{Round: "R32", SlotCode: "R32-04", Label: ptr("Round of 32 - Match 4"), HomeLabel: ptr("1D"), AwayLabel: ptr("3A/B/C"), Kickoff: parseDate("2026-06-29 20:00:00")},
		models.KnockoutSlot{Round: "R32", SlotCode: "R32-05", Label: ptr("Round of 32 - Match 5"), HomeLabel: ptr("1E"), AwayLabel: ptr("2F"), Kickoff: parseDate("2026-06-30 16:00:00")},
		models.KnockoutSlot{Round: "R32", SlotCode: "R32-06", Label: ptr("Round of 32 - Match 6"), HomeLabel: ptr("1F"), AwayLabel: ptr("2E"), Kickoff: parseDate("2026-06-30 20:00:00")},
		models.KnockoutSlot{Round: "R32", SlotCode: "R32-07", Label: ptr("Round of 32 - Match 7"), HomeLabel: ptr("1G"), AwayLabel: ptr("2H"), Kickoff: parseDate("2026-07-01 16:00:00")},
		models.KnockoutSlot{Round: "R32", SlotCode: "R32-08", Label: ptr("Round of 32 - Match 8"), HomeLabel: ptr("1H"), AwayLabel: ptr("2G"), Kickoff: parseDate("2026-07-01 20:00:00")},
		models.KnockoutSlot{Round: "R32", SlotCode: "R32-09", Label: ptr("Round of 32 - Match 9"), HomeLabel: ptr("1I"), AwayLabel: ptr("2J"), Kickoff: parseDate("2026-07-02 16:00:00")},
		models.KnockoutSlot{Round: "R32", SlotCode: "R32-10", Label: ptr("Round of 32 - Match 10"), HomeLabel: ptr("1J"), AwayLabel: ptr("2I"), Kickoff: parseDate("2026-07-02 20:00:00")},
		models.KnockoutSlot{Round: "R32", SlotCode: "R32-11", Label: ptr("Round of 32 - Match 11"), HomeLabel: ptr("1K"), AwayLabel: ptr("2L"), Kickoff: parseDate("2026-07-03 16:00:00")},
		models.KnockoutSlot{Round: "R32", SlotCode: "R32-12", Label: ptr("Round of 32 - Match 12"), HomeLabel: ptr("1L"), AwayLabel: ptr("2K"), Kickoff: parseDate("2026-07-03 20:00:00")},
		models.KnockoutSlot{Round: "R32", SlotCode: "R32-13", Label: ptr("Round of 32 - Match 13"), HomeLabel: ptr("2A"), AwayLabel: ptr("2B"), Kickoff: parseDate("2026-07-04 16:00:00")},
		models.KnockoutSlot{Round: "R32", SlotCode: "R32-14", Label: ptr("Round of 32 - Match 14"), HomeLabel: ptr("2C"), AwayLabel: ptr("2D"), Kickoff: parseDate("2026-07-04 20:00:00")},
		models.KnockoutSlot{Round: "R32", SlotCode: "R32-15", Label: ptr("Round of 32 - Match 15"), HomeLabel: ptr("3F/G/H"), AwayLabel: ptr("3I/J/K"), Kickoff: parseDate("2026-07-05 16:00:00")},
		models.KnockoutSlot{Round: "R32", SlotCode: "R32-16", Label: ptr("Round of 32 - Match 16"), HomeLabel: ptr("3L/A/B"), AwayLabel: ptr("3C/D/E"), Kickoff: parseDate("2026-07-05 20:00:00")},
		models.KnockoutSlot{Round: "R16", SlotCode: "R16-01", Label: ptr("Round of 16 - Match 1"), HomeLabel: ptr("W R32-01"), AwayLabel: ptr("W R32-02"), Kickoff: parseDate("2026-07-07 16:00:00")},
		models.KnockoutSlot{Round: "R16", SlotCode: "R16-02", Label: ptr("Round of 16 - Match 2"), HomeLabel: ptr("W R32-03"), AwayLabel: ptr("W R32-04"), Kickoff: parseDate("2026-07-07 20:00:00")},
		models.KnockoutSlot{Round: "R16", SlotCode: "R16-03", Label: ptr("Round of 16 - Match 3"), HomeLabel: ptr("W R32-05"), AwayLabel: ptr("W R32-06"), Kickoff: parseDate("2026-07-08 16:00:00")},
		models.KnockoutSlot{Round: "R16", SlotCode: "R16-04", Label: ptr("Round of 16 - Match 4"), HomeLabel: ptr("W R32-07"), AwayLabel: ptr("W R32-08"), Kickoff: parseDate("2026-07-08 20:00:00")},
		models.KnockoutSlot{Round: "R16", SlotCode: "R16-05", Label: ptr("Round of 16 - Match 5"), HomeLabel: ptr("W R32-09"), AwayLabel: ptr("W R32-10"), Kickoff: parseDate("2026-07-09 16:00:00")},
		models.KnockoutSlot{Round: "R16", SlotCode: "R16-06", Label: ptr("Round of 16 - Match 6"), HomeLabel: ptr("W R32-11"), AwayLabel: ptr("W R32-12"), Kickoff: parseDate("2026-07-09 20:00:00")},
		models.KnockoutSlot{Round: "R16", SlotCode: "R16-07", Label: ptr("Round of 16 - Match 7"), HomeLabel: ptr("W R32-13"), AwayLabel: ptr("W R32-14"), Kickoff: parseDate("2026-07-10 16:00:00")},
		models.KnockoutSlot{Round: "R16", SlotCode: "R16-08", Label: ptr("Round of 16 - Match 8"), HomeLabel: ptr("W R32-15"), AwayLabel: ptr("W R32-16"), Kickoff: parseDate("2026-07-10 20:00:00")},
		models.KnockoutSlot{Round: "QF", SlotCode: "QF-01", Label: ptr("Quarterfinal 1"), HomeLabel: ptr("W R16-01"), AwayLabel: ptr("W R16-02"), Kickoff: parseDate("2026-07-11 16:00:00")},
		models.KnockoutSlot{Round: "QF", SlotCode: "QF-02", Label: ptr("Quarterfinal 2"), HomeLabel: ptr("W R16-03"), AwayLabel: ptr("W R16-04"), Kickoff: parseDate("2026-07-11 20:00:00")},
		models.KnockoutSlot{Round: "QF", SlotCode: "QF-03", Label: ptr("Quarterfinal 3"), HomeLabel: ptr("W R16-05"), AwayLabel: ptr("W R16-06"), Kickoff: parseDate("2026-07-12 16:00:00")},
		models.KnockoutSlot{Round: "QF", SlotCode: "QF-04", Label: ptr("Quarterfinal 4"), HomeLabel: ptr("W R16-07"), AwayLabel: ptr("W R16-08"), Kickoff: parseDate("2026-07-12 20:00:00")},
		models.KnockoutSlot{Round: "SF", SlotCode: "SF-01", Label: ptr("Semifinal 1"), HomeLabel: ptr("W QF-01"), AwayLabel: ptr("W QF-02"), Kickoff: parseDate("2026-07-14 20:00:00")},
		models.KnockoutSlot{Round: "SF", SlotCode: "SF-02", Label: ptr("Semifinal 2"), HomeLabel: ptr("W QF-03"), AwayLabel: ptr("W QF-04"), Kickoff: parseDate("2026-07-15 20:00:00")},
		models.KnockoutSlot{Round: "3RD", SlotCode: "3RD-01", Label: ptr("Third Place Match"), HomeLabel: ptr("L SF-01"), AwayLabel: ptr("L SF-02"), Kickoff: parseDate("2026-07-18 16:00:00")},
		models.KnockoutSlot{Round: "FINAL", SlotCode: "FINAL", Label: ptr("Final"), HomeLabel: ptr("W SF-01"), AwayLabel: ptr("W SF-02"), Kickoff: parseDate("2026-07-19 20:00:00")},
	}
	_, err = slotsCol.InsertMany(ctx, slots)
	if err != nil {
		log.Printf("Failed to seed knockout slots: %v", err)
	}

	// 4. Matches
	matchesCol := db.Collection("matches")
	matchesCol.Drop(ctx)
	matches := []models.Match{
		{GroupCode: "A", MatchNumber: ptrInt(1), HomeTeam: "MEX", AwayTeam: "RSA", Kickoff: *parseDate("2026-06-11 20:00:00")},
		{GroupCode: "A", MatchNumber: ptrInt(2), HomeTeam: "KOR", AwayTeam: "CZE", Kickoff: *parseDate("2026-06-12 02:00:00")},
		{GroupCode: "A", MatchNumber: ptrInt(25), HomeTeam: "MEX", AwayTeam: "KOR", Kickoff: *parseDate("2026-06-16 20:00:00")},
		{GroupCode: "A", MatchNumber: ptrInt(26), HomeTeam: "RSA", AwayTeam: "CZE", Kickoff: *parseDate("2026-06-17 01:00:00")},
		{GroupCode: "A", MatchNumber: ptrInt(49), HomeTeam: "RSA", AwayTeam: "KOR", Kickoff: *parseDate("2026-06-22 18:00:00")},
		{GroupCode: "A", MatchNumber: ptrInt(50), HomeTeam: "MEX", AwayTeam: "CZE", Kickoff: *parseDate("2026-06-22 18:00:00")},
		{GroupCode: "B", MatchNumber: ptrInt(3), HomeTeam: "CAN", AwayTeam: "BIH", Kickoff: *parseDate("2026-06-12 19:00:00")},
		{GroupCode: "B", MatchNumber: ptrInt(12), HomeTeam: "QAT", AwayTeam: "SUI", Kickoff: *parseDate("2026-06-14 03:00:00")},
		{GroupCode: "B", MatchNumber: ptrInt(27), HomeTeam: "CAN", AwayTeam: "QAT", Kickoff: *parseDate("2026-06-17 19:00:00")},
		{GroupCode: "B", MatchNumber: ptrInt(28), HomeTeam: "BIH", AwayTeam: "SUI", Kickoff: *parseDate("2026-06-18 01:00:00")},
		{GroupCode: "B", MatchNumber: ptrInt(51), HomeTeam: "BIH", AwayTeam: "QAT", Kickoff: *parseDate("2026-06-23 18:00:00")},
		{GroupCode: "B", MatchNumber: ptrInt(52), HomeTeam: "CAN", AwayTeam: "SUI", Kickoff: *parseDate("2026-06-23 18:00:00")},
		{GroupCode: "C", MatchNumber: ptrInt(5), HomeTeam: "BRA", AwayTeam: "MAR", Kickoff: *parseDate("2026-06-13 16:00:00")},
		{GroupCode: "C", MatchNumber: ptrInt(6), HomeTeam: "HAI", AwayTeam: "SCO", Kickoff: *parseDate("2026-06-13 19:00:00")},
		{GroupCode: "C", MatchNumber: ptrInt(29), HomeTeam: "BRA", AwayTeam: "HAI", Kickoff: *parseDate("2026-06-18 18:00:00")},
		{GroupCode: "C", MatchNumber: ptrInt(30), HomeTeam: "MAR", AwayTeam: "SCO", Kickoff: *parseDate("2026-06-19 01:00:00")},
		{GroupCode: "C", MatchNumber: ptrInt(53), HomeTeam: "MAR", AwayTeam: "HAI", Kickoff: *parseDate("2026-06-24 18:00:00")},
		{GroupCode: "C", MatchNumber: ptrInt(54), HomeTeam: "BRA", AwayTeam: "SCO", Kickoff: *parseDate("2026-06-24 18:00:00")},
		{GroupCode: "D", MatchNumber: ptrInt(4), HomeTeam: "USA", AwayTeam: "PAR", Kickoff: *parseDate("2026-06-13 02:00:00")},
		{GroupCode: "D", MatchNumber: ptrInt(7), HomeTeam: "AUS", AwayTeam: "TUR", Kickoff: *parseDate("2026-06-14 02:00:00")},
		{GroupCode: "D", MatchNumber: ptrInt(31), HomeTeam: "USA", AwayTeam: "AUS", Kickoff: *parseDate("2026-06-18 22:00:00")},
		{GroupCode: "D", MatchNumber: ptrInt(32), HomeTeam: "PAR", AwayTeam: "TUR", Kickoff: *parseDate("2026-06-19 03:00:00")},
		{GroupCode: "D", MatchNumber: ptrInt(55), HomeTeam: "PAR", AwayTeam: "AUS", Kickoff: *parseDate("2026-06-25 02:00:00")},
		{GroupCode: "D", MatchNumber: ptrInt(56), HomeTeam: "USA", AwayTeam: "TUR", Kickoff: *parseDate("2026-06-25 02:00:00")},
		{GroupCode: "E", MatchNumber: ptrInt(8), HomeTeam: "GER", AwayTeam: "CUW", Kickoff: *parseDate("2026-06-14 16:00:00")},
		{GroupCode: "E", MatchNumber: ptrInt(9), HomeTeam: "CIV", AwayTeam: "ECU", Kickoff: *parseDate("2026-06-14 22:00:00")},
		{GroupCode: "E", MatchNumber: ptrInt(33), HomeTeam: "GER", AwayTeam: "CIV", Kickoff: *parseDate("2026-06-20 19:00:00")},
		{GroupCode: "E", MatchNumber: ptrInt(34), HomeTeam: "ECU", AwayTeam: "CUW", Kickoff: *parseDate("2026-06-20 22:00:00")},
		{GroupCode: "E", MatchNumber: ptrInt(57), HomeTeam: "CUW", AwayTeam: "CIV", Kickoff: *parseDate("2026-06-25 19:00:00")},
		{GroupCode: "E", MatchNumber: ptrInt(58), HomeTeam: "ECU", AwayTeam: "GER", Kickoff: *parseDate("2026-06-25 19:00:00")},
		{GroupCode: "F", MatchNumber: ptrInt(10), HomeTeam: "NED", AwayTeam: "JPN", Kickoff: *parseDate("2026-06-15 00:00:00")},
		{GroupCode: "F", MatchNumber: ptrInt(11), HomeTeam: "SWE", AwayTeam: "TUN", Kickoff: *parseDate("2026-06-15 03:00:00")},
		{GroupCode: "F", MatchNumber: ptrInt(35), HomeTeam: "NED", AwayTeam: "SWE", Kickoff: *parseDate("2026-06-20 16:00:00")},
		{GroupCode: "F", MatchNumber: ptrInt(36), HomeTeam: "JPN", AwayTeam: "TUN", Kickoff: *parseDate("2026-06-21 01:00:00")},
		{GroupCode: "F", MatchNumber: ptrInt(59), HomeTeam: "JPN", AwayTeam: "SWE", Kickoff: *parseDate("2026-06-25 22:00:00")},
		{GroupCode: "F", MatchNumber: ptrInt(60), HomeTeam: "TUN", AwayTeam: "NED", Kickoff: *parseDate("2026-06-25 22:00:00")},
		{GroupCode: "G", MatchNumber: ptrInt(15), HomeTeam: "BEL", AwayTeam: "EGY", Kickoff: *parseDate("2026-06-16 02:00:00")},
		{GroupCode: "G", MatchNumber: ptrInt(16), HomeTeam: "IRN", AwayTeam: "NZL", Kickoff: *parseDate("2026-06-16 05:00:00")},
		{GroupCode: "G", MatchNumber: ptrInt(37), HomeTeam: "BEL", AwayTeam: "IRN", Kickoff: *parseDate("2026-06-21 02:00:00")},
		{GroupCode: "G", MatchNumber: ptrInt(38), HomeTeam: "NZL", AwayTeam: "EGY", Kickoff: *parseDate("2026-06-21 19:00:00")},
		{GroupCode: "G", MatchNumber: ptrInt(61), HomeTeam: "EGY", AwayTeam: "IRN", Kickoff: *parseDate("2026-06-26 02:00:00")},
		{GroupCode: "G", MatchNumber: ptrInt(62), HomeTeam: "NZL", AwayTeam: "BEL", Kickoff: *parseDate("2026-06-26 02:00:00")},
		{GroupCode: "H", MatchNumber: ptrInt(13), HomeTeam: "ESP", AwayTeam: "CPV", Kickoff: *parseDate("2026-06-15 19:00:00")},
		{GroupCode: "H", MatchNumber: ptrInt(14), HomeTeam: "KSA", AwayTeam: "URU", Kickoff: *parseDate("2026-06-15 22:00:00")},
		{GroupCode: "H", MatchNumber: ptrInt(39), HomeTeam: "ESP", AwayTeam: "KSA", Kickoff: *parseDate("2026-06-21 16:00:00")},
		{GroupCode: "H", MatchNumber: ptrInt(40), HomeTeam: "URU", AwayTeam: "CPV", Kickoff: *parseDate("2026-06-21 22:00:00")},
		{GroupCode: "H", MatchNumber: ptrInt(63), HomeTeam: "CPV", AwayTeam: "KSA", Kickoff: *parseDate("2026-06-26 19:00:00")},
		{GroupCode: "H", MatchNumber: ptrInt(64), HomeTeam: "URU", AwayTeam: "ESP", Kickoff: *parseDate("2026-06-26 19:00:00")},
		{GroupCode: "I", MatchNumber: ptrInt(17), HomeTeam: "FRA", AwayTeam: "SEN", Kickoff: *parseDate("2026-06-16 19:00:00")},
		{GroupCode: "I", MatchNumber: ptrInt(18), HomeTeam: "IRQ", AwayTeam: "NOR", Kickoff: *parseDate("2026-06-16 22:00:00")},
		{GroupCode: "I", MatchNumber: ptrInt(41), HomeTeam: "FRA", AwayTeam: "IRQ", Kickoff: *parseDate("2026-06-22 19:00:00")},
		{GroupCode: "I", MatchNumber: ptrInt(42), HomeTeam: "NOR", AwayTeam: "SEN", Kickoff: *parseDate("2026-06-22 22:00:00")},
		{GroupCode: "I", MatchNumber: ptrInt(65), HomeTeam: "NOR", AwayTeam: "FRA", Kickoff: *parseDate("2026-06-26 16:00:00")},
		{GroupCode: "I", MatchNumber: ptrInt(66), HomeTeam: "SEN", AwayTeam: "IRQ", Kickoff: *parseDate("2026-06-26 16:00:00")},
		{GroupCode: "J", MatchNumber: ptrInt(19), HomeTeam: "ARG", AwayTeam: "ALG", Kickoff: *parseDate("2026-06-17 00:00:00")},
		{GroupCode: "J", MatchNumber: ptrInt(20), HomeTeam: "AUT", AwayTeam: "JOR", Kickoff: *parseDate("2026-06-17 03:00:00")},
		{GroupCode: "J", MatchNumber: ptrInt(43), HomeTeam: "ARG", AwayTeam: "AUT", Kickoff: *parseDate("2026-06-22 02:00:00")},
		{GroupCode: "J", MatchNumber: ptrInt(44), HomeTeam: "JOR", AwayTeam: "ALG", Kickoff: *parseDate("2026-06-22 05:00:00")},
		{GroupCode: "J", MatchNumber: ptrInt(67), HomeTeam: "ALG", AwayTeam: "AUT", Kickoff: *parseDate("2026-06-27 00:00:00")},
		{GroupCode: "J", MatchNumber: ptrInt(68), HomeTeam: "JOR", AwayTeam: "ARG", Kickoff: *parseDate("2026-06-27 00:00:00")},
		{GroupCode: "K", MatchNumber: ptrInt(21), HomeTeam: "POR", AwayTeam: "COD", Kickoff: *parseDate("2026-06-17 16:00:00")},
		{GroupCode: "K", MatchNumber: ptrInt(22), HomeTeam: "UZB", AwayTeam: "COL", Kickoff: *parseDate("2026-06-17 19:00:00")},
		{GroupCode: "K", MatchNumber: ptrInt(45), HomeTeam: "POR", AwayTeam: "UZB", Kickoff: *parseDate("2026-06-23 16:00:00")},
		{GroupCode: "K", MatchNumber: ptrInt(46), HomeTeam: "COL", AwayTeam: "COD", Kickoff: *parseDate("2026-06-23 19:00:00")},
		{GroupCode: "K", MatchNumber: ptrInt(69), HomeTeam: "COL", AwayTeam: "POR", Kickoff: *parseDate("2026-06-27 19:00:00")},
		{GroupCode: "K", MatchNumber: ptrInt(70), HomeTeam: "COD", AwayTeam: "UZB", Kickoff: *parseDate("2026-06-27 19:00:00")},
		{GroupCode: "L", MatchNumber: ptrInt(23), HomeTeam: "ENG", AwayTeam: "CRO", Kickoff: *parseDate("2026-06-18 00:00:00")},
		{GroupCode: "L", MatchNumber: ptrInt(24), HomeTeam: "GHA", AwayTeam: "PAN", Kickoff: *parseDate("2026-06-18 03:00:00")},
		{GroupCode: "L", MatchNumber: ptrInt(47), HomeTeam: "ENG", AwayTeam: "GHA", Kickoff: *parseDate("2026-06-23 22:00:00")},
		{GroupCode: "L", MatchNumber: ptrInt(48), HomeTeam: "PAN", AwayTeam: "CRO", Kickoff: *parseDate("2026-06-24 01:00:00")},
		{GroupCode: "L", MatchNumber: ptrInt(71), HomeTeam: "PAN", AwayTeam: "ENG", Kickoff: *parseDate("2026-06-28 00:00:00")},
		{GroupCode: "L", MatchNumber: ptrInt(72), HomeTeam: "CRO", AwayTeam: "GHA", Kickoff: *parseDate("2026-06-28 00:00:00")},
	}
	interfaceMatches := make([]interface{}, len(matches))
	for i, m := range matches {
		interfaceMatches[i] = m
	}
	_, err = matchesCol.InsertMany(ctx, interfaceMatches)
	if err != nil {
		log.Printf("Failed to seed matches: %v", err)
	}

	log.Println("Seeding completed successfully")
}

func ptr(s string) *string { return &s }
func ptrInt(i int) *int    { return &i }
func parseDate(s string) *time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", s)
	return &t
}
