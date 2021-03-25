package purdoobahs

type IPurdoobahService interface {
	All() ([]*Purdoobah, error)
	ByName(string) (*Purdoobah, error)
}

type Purdoobah struct {
	Name                 string `json:"name"`
	BirthCertificateName struct {
		First  string `json:"first"`
		Middle string `json:"middle,omitempty"`
		Last   string `json:"last"`
	} `json:"birth_certificate_name"`
	Emoji string `json:"emoji"`

	Marching struct {
		YearsMarched []int  `json:"years_marched"`
		Shoutout     string `json:"shoutout,omitempty"`
	} `json:"marching"`

	Education struct {
		Major string `json:"major"`
		Minor string `json:"minor,omitempty"`
		Year  string `json:"year"`
	} `json:"education"`

	Hometown struct {
		City  string `json:"city"`
		State string `json:"state"`
	} `json:"hometown"`

	Alumni struct {
		Job string `json:"job,omitempty"`
	} `json:"alumni,omitempty"`

	Personal struct {
		Hobbies []string `json:"hobbies,omitempty"`
		Socials struct {
			Facebook  string `json:"facebook,omitempty"`
			Instagram string `json:"instagram,omitempty"`
			LinkedIn  string `json:"linkedin,omitempty"`
		} `json:"socials,omitempty"`
	} `json:"personal,omitempty"`

	Achievements struct {
		StudentLeader         []int `json:"student_leader,omitempty"`
		BottomFeederCommittee bool  `json:"bottom_feeder_committee,omitempty"`
		SpoonsassinsVictories []int `json:"spoonsassins_victories,omitempty"`
		KappaKappaPsi         bool  `json:"kappa_kappa_psi,omitempty"`
		TauBetaSigma          bool  `json:"tau_beta_sigma,omitempty"`
	} `json:"achievements,omitempty"`

	Metadata struct {
		Image string `json:"image"`
	} `json:"metadata"`
}
