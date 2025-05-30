package db

import (
	"context"
	"log"
	"math/rand"

	"faizisyellow.github.com/thegosocialnetwork/internal/store"
	"github.com/brianvoe/gofakeit/v7"
)

var titles = []string{
	"Wandering Through Silence",
	"Chasing Starlight",
	"Beneath the Willow",
	"Moments Between Raindrops",
	"Painted Horizons",
	"Shadows of Yesterday",
	"When the Earth Paused",
	"Footprints in Fog",
	"Lanterns in the Wind",
	"The Sound of Empty",
	"Between Now and Never",
	"Paper Hearts",
	"Letters Never Sent",
	"Skylines and Skylights",
	"Echoes of Laughter",
	"The Art of Waiting",
	"Whispers from the Attic",
	"Of Coffee and Chaos",
	"Fractured Light",
	"Melodies Unwritten",
	"Cracked Teacups",
	"Seasons We Forgot",
	"Maps Without Roads",
	"The Room With No Clock",
	"Rain on Old Rooftops",
	"Blank Pages",
	"Clouds Like Canvas",
	"The Window Seat",
	"Fiction in Reality",
	"Lanterns on Water",
	"Rust and Romance",
	"Secondhand Stars",
	"Curtains Half Drawn",
	"Shoes by the Door",
	"Hands Without Rings",
	"Tides and Time",
	"Windswept Letters",
	"Hearts Like Keys",
	"Photographs that Breathe",
	"Ink-Stained Smiles",
	"The Quiet Kind",
	"Light Behind Doors",
	"The Weight of Small Things",
	"Fire Escapes",
	"Rain Checks and Regrets",
	"Daydream Drafts",
	"Walls Can Listen",
	"Books with No Titles",
	"The Clock Skipped",
	"Unfolding Maps",
	"Wires and Wishes",
}

var contents = []string{
	"In the hush of dawn, every whisper carries a story untold.",
	"We lose ourselves in galaxies that remind us how small and infinite we are.",
	"The old tree knows secrets that even time has forgotten.",
	"Sometimes, peace is found not in the storm, but in the stillness between it.",
	"Every sunset bleeds a new memory into the sky.",
	"We walk forward, but shadows of our past never quite leave us.",
	"There are moments where everything halts, and that silence changes everything.",
	"Some paths we follow without knowing where they lead, only that we must.",
	"Hope flickers in the gentlest winds, yet it persists.",
	"Not all silence is the same; some is loud in memory.",
	"Time isn’t always a line—it bends, folds, and loops with emotion.",
	"Fragile yet intentional, the way we give our hearts matters most.",
	"Words once bottled still echo in the corridors of time.",
	"The city breathes light even in its darkest corners.",
	"Some laughs haunt us sweetly, like ghosts that mean no harm.",
	"There is beauty in patience, in the slow unfolding of things.",
	"Dust holds more memories than we give it credit for.",
	"Life brews better when it's a little bitter, a little bold.",
	"Even broken glass catches the sun if you hold it right.",
	"The best songs are the ones we feel but never play.",
	"Some things aren’t ruined, just more interesting with their flaws.",
	"Not all years are counted in calendars.",
	"Wanderers need direction, not highways.",
	"Time feels different where we stop measuring it.",
	"There’s a comfort in hearing what the past still sounds like.",
	"Not all stories start with ink.",
	"Even the sky tries to be art sometimes.",
	"Where we sit often changes what we see.",
	"The world writes plots stranger than novels.",
	"Light dances differently when it floats.",
	"Even decay can be beautiful when remembered right.",
	"We all shine a little from someone else's spark.",
	"Some truths hide better in the almost-seen.",
	"Homes begin where footsteps end.",
	"Commitment isn’t always a circle.",
	"What the sea takes, it sometimes returns.",
	"Paper can travel when it knows what to say.",
	"We’re all trying to unlock something in someone.",
	"Still images move something inside us.",
	"Writers leave parts of themselves on every page.",
	"Not all strength roars.",
	"Hope glows brightest when hidden.",
	"Tiny details make heavy memories.",
	"Some exits lead to the best views.",
	"Plans change, but feelings linger.",
	"Imagination is just reality with the volume turned down.",
	"Even silence hears us.",
	"Not everything needs a label to matter.",
	"Some hours aren’t lost—they’re stolen.",
	"Journeys start before we know it.",
	"Connection isn’t always visible, but always felt.",
}

func Seed(store store.Storage) {
	ctx := context.Background()

	// users := generateUsers(10)
	// for _, user := range users {
	// 	store.Users.Create(ctx, user)
	// }

	posts := generatePost(100)

	for _, post := range posts {
		store.Posts.Create(ctx, post)
	}

	log.Println("Seeding successfully...")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	ps := gofakeit.Password(true, false, false, true, false, 5)

	for i, _ := range users {
		users[i] = &store.User{
			Username: gofakeit.Username(),
			Email:    gofakeit.Email(),
			Password: store.HashPassword{
				Text: &ps,
				Hash: []byte("helloworld"),
			},
		}
	}

	return users
}

func generatePost(num int) []*store.Post {
	posts := make([]*store.Post, num)

	for i, _ := range posts {
		// length of rows users table hehehe...
		iduser := rand.Intn(380)

		// id users below 4 already removed
		if iduser >= 0 && iduser <= 4 {
			iduser *= 6
		}

		posts[i] = &store.Post{
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			UserID:  iduser,
		}

	}

	return posts
}
