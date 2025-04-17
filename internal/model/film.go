package model

type FilmInfo struct {
    Title           string
    Tagline         string
    Synopsis        string
    Director        string
    Producer        string
    ContactEmail    string
    ContactPhone    string
    DurationMinutes int
    Genre           string
    Locations       []Location
    TeamMembers     []TeamMember
    Cast            []CastMember
}

type Location struct {
    Name        string
    Description string
    ImageURL    string
}

type TeamMember struct {
    Role     string
    Name     string
    Email    string
    ImageURL string
}

type CastMember struct {
    Role      string
    ActorName string
    ImageURL  string
}