package nba

import "time"

type NBATeam struct {
	Sports []struct {
		ID    string `json:"id"`
		UID   string `json:"uid"`
		GUID  string `json:"guid"`
		Name  string `json:"name"`
		Slug  string `json:"slug"`
		Logos []struct {
			Href   string   `json:"href"`
			Alt    string   `json:"alt"`
			Rel    []string `json:"rel"`
			Width  int      `json:"width"`
			Height int      `json:"height"`
		} `json:"logos"`
		Leagues []struct {
			ID           string   `json:"id"`
			UID          string   `json:"uid"`
			Name         string   `json:"name"`
			Abbreviation string   `json:"abbreviation"`
			ShortName    string   `json:"shortName"`
			Slug         string   `json:"slug"`
			Tag          string   `json:"tag"`
			IsTournament bool     `json:"isTournament"`
			Smartdates   []string `json:"smartdates"`
			Events       []struct {
				ID            string    `json:"id"`
				UID           string    `json:"uid"`
				Date          time.Time `json:"date"`
				TimeValid     bool      `json:"timeValid"`
				Recent        bool      `json:"recent"`
				Name          string    `json:"name"`
				ShortName     string    `json:"shortName"`
				SeriesSummary string    `json:"seriesSummary"`
				Links         []struct {
					Rel  []string `json:"rel"`
					Href string   `json:"href"`
					Text string   `json:"text"`
				} `json:"links"`
				GamecastAvailable   bool      `json:"gamecastAvailable"`
				PlayByPlayAvailable bool      `json:"playByPlayAvailable"`
				CommentaryAvailable bool      `json:"commentaryAvailable"`
				OnWatch             bool      `json:"onWatch"`
				CompetitionID       string    `json:"competitionId"`
				Location            string    `json:"location"`
				Season              int       `json:"season"`
				SeasonStartDate     time.Time `json:"seasonStartDate"`
				SeasonEndDate       time.Time `json:"seasonEndDate"`
				SeasonType          string    `json:"seasonType"`
				SeasonTypeHasGroups bool      `json:"seasonTypeHasGroups"`
				Group               struct {
					GroupID      string `json:"groupId"`
					Name         string `json:"name"`
					Abbreviation string `json:"abbreviation"`
					ShortName    string `json:"shortName"`
				} `json:"group"`
				Week       int    `json:"week"`
				WeekText   string `json:"weekText"`
				Link       string `json:"link"`
				Status     string `json:"status"`
				Summary    string `json:"summary"`
				Period     int    `json:"period"`
				Clock      string `json:"clock"`
				Broadcasts []struct {
					TypeID        int    `json:"typeId"`
					Priority      int    `json:"priority"`
					Type          string `json:"type"`
					IsNational    bool   `json:"isNational"`
					BroadcasterID int    `json:"broadcasterId"`
					BroadcastID   int    `json:"broadcastId"`
					Name          string `json:"name"`
					ShortName     string `json:"shortName"`
					CallLetters   string `json:"callLetters"`
					Station       string `json:"station"`
					Lang          string `json:"lang"`
					Region        string `json:"region"`
					Slug          string `json:"slug"`
				} `json:"broadcasts,omitempty"`
				Broadcast string `json:"broadcast,omitempty"`
				Odds      struct {
					Details   string  `json:"details"`
					OverUnder float64 `json:"overUnder"`
					Spread    float64 `json:"spread"`
					OverOdds  int     `json:"overOdds"`
					UnderOdds int     `json:"underOdds"`
					Provider  struct {
						ID       string `json:"id"`
						Name     string `json:"name"`
						Priority int    `json:"priority"`
						Logos    []struct {
							Href string   `json:"href"`
							Rel  []string `json:"rel"`
						} `json:"logos"`
					} `json:"provider"`
					Home struct {
						MoneyLine int `json:"moneyLine"`
					} `json:"home"`
					Away struct {
						MoneyLine int `json:"moneyLine"`
					} `json:"away"`
					AwayTeamOdds struct {
						Favorite   bool `json:"favorite"`
						Underdog   bool `json:"underdog"`
						MoneyLine  int  `json:"moneyLine"`
						SpreadOdds int  `json:"spreadOdds"`
						Team       struct {
							ID           string `json:"id"`
							Abbreviation string `json:"abbreviation"`
						} `json:"team"`
					} `json:"awayTeamOdds"`
					HomeTeamOdds struct {
						Favorite   bool `json:"favorite"`
						Underdog   bool `json:"underdog"`
						MoneyLine  int  `json:"moneyLine"`
						SpreadOdds int  `json:"spreadOdds"`
						Team       struct {
							ID           string `json:"id"`
							Abbreviation string `json:"abbreviation"`
						} `json:"team"`
					} `json:"homeTeamOdds"`
					Links []struct {
						Language   string   `json:"language"`
						Rel        []string `json:"rel"`
						Href       string   `json:"href"`
						Text       string   `json:"text"`
						ShortText  string   `json:"shortText"`
						IsExternal bool     `json:"isExternal"`
						IsPremium  bool     `json:"isPremium"`
					} `json:"links"`
					PointSpread struct {
						DisplayName      string `json:"displayName"`
						ShortDisplayName string `json:"shortDisplayName"`
						Home             struct {
							Open struct {
								Line string `json:"line"`
								Odds string `json:"odds"`
								Link struct {
									Rel      []string `json:"rel"`
									Href     string   `json:"href"`
									Tracking struct {
										Campaign string `json:"campaign"`
										Tags     struct {
											League     string `json:"league"`
											Sport      string `json:"sport"`
											GameID     int    `json:"gameId"`
											BetSide    string `json:"betSide"`
											BetType    string `json:"betType"`
											BetDetails string `json:"betDetails"`
										} `json:"tags"`
									} `json:"tracking"`
								} `json:"link"`
							} `json:"open"`
							Close struct {
								Line string `json:"line"`
								Odds string `json:"odds"`
								Link struct {
									Rel      []string `json:"rel"`
									Href     string   `json:"href"`
									Tracking struct {
										Campaign string `json:"campaign"`
										Tags     struct {
											League     string `json:"league"`
											Sport      string `json:"sport"`
											GameID     int    `json:"gameId"`
											BetSide    string `json:"betSide"`
											BetType    string `json:"betType"`
											BetDetails string `json:"betDetails"`
										} `json:"tags"`
									} `json:"tracking"`
								} `json:"link"`
							} `json:"close"`
						} `json:"home"`
						Away struct {
							Open struct {
								Line string `json:"line"`
								Odds string `json:"odds"`
								Link struct {
									Rel      []string `json:"rel"`
									Href     string   `json:"href"`
									Tracking struct {
										Campaign string `json:"campaign"`
										Tags     struct {
											League     string `json:"league"`
											Sport      string `json:"sport"`
											GameID     int    `json:"gameId"`
											BetSide    string `json:"betSide"`
											BetType    string `json:"betType"`
											BetDetails string `json:"betDetails"`
										} `json:"tags"`
									} `json:"tracking"`
								} `json:"link"`
							} `json:"open"`
							Close struct {
								Line string `json:"line"`
								Odds string `json:"odds"`
								Link struct {
									Rel      []string `json:"rel"`
									Href     string   `json:"href"`
									Tracking struct {
										Campaign string `json:"campaign"`
										Tags     struct {
											League     string `json:"league"`
											Sport      string `json:"sport"`
											GameID     int    `json:"gameId"`
											BetSide    string `json:"betSide"`
											BetType    string `json:"betType"`
											BetDetails string `json:"betDetails"`
										} `json:"tags"`
									} `json:"tracking"`
								} `json:"link"`
							} `json:"close"`
						} `json:"away"`
					} `json:"pointSpread"`
					Moneyline struct {
						DisplayName      string `json:"displayName"`
						ShortDisplayName string `json:"shortDisplayName"`
						Home             struct {
							Open struct {
								Odds string `json:"odds"`
								Link struct {
									Rel      []string `json:"rel"`
									Href     string   `json:"href"`
									Tracking struct {
										Campaign string `json:"campaign"`
										Tags     struct {
											League  string `json:"league"`
											Sport   string `json:"sport"`
											GameID  int    `json:"gameId"`
											BetSide string `json:"betSide"`
											BetType string `json:"betType"`
										} `json:"tags"`
									} `json:"tracking"`
								} `json:"link"`
							} `json:"open"`
							Close struct {
								Odds string `json:"odds"`
								Link struct {
									Rel      []string `json:"rel"`
									Href     string   `json:"href"`
									Tracking struct {
										Campaign string `json:"campaign"`
										Tags     struct {
											League  string `json:"league"`
											Sport   string `json:"sport"`
											GameID  int    `json:"gameId"`
											BetSide string `json:"betSide"`
											BetType string `json:"betType"`
										} `json:"tags"`
									} `json:"tracking"`
								} `json:"link"`
							} `json:"close"`
						} `json:"home"`
						Away struct {
							Open struct {
								Odds string `json:"odds"`
								Link struct {
									Rel      []string `json:"rel"`
									Href     string   `json:"href"`
									Tracking struct {
										Campaign string `json:"campaign"`
										Tags     struct {
											League  string `json:"league"`
											Sport   string `json:"sport"`
											GameID  int    `json:"gameId"`
											BetSide string `json:"betSide"`
											BetType string `json:"betType"`
										} `json:"tags"`
									} `json:"tracking"`
								} `json:"link"`
							} `json:"open"`
							Close struct {
								Odds string `json:"odds"`
								Link struct {
									Rel      []string `json:"rel"`
									Href     string   `json:"href"`
									Tracking struct {
										Campaign string `json:"campaign"`
										Tags     struct {
											League  string `json:"league"`
											Sport   string `json:"sport"`
											GameID  int    `json:"gameId"`
											BetSide string `json:"betSide"`
											BetType string `json:"betType"`
										} `json:"tags"`
									} `json:"tracking"`
								} `json:"link"`
							} `json:"close"`
						} `json:"away"`
					} `json:"moneyline"`
					Total struct {
						DisplayName      string `json:"displayName"`
						ShortDisplayName string `json:"shortDisplayName"`
						Over             struct {
							Open struct {
								Line string `json:"line"`
								Odds string `json:"odds"`
								Link struct {
									Rel      []string `json:"rel"`
									Href     string   `json:"href"`
									Tracking struct {
										Campaign string `json:"campaign"`
										Tags     struct {
											League     string `json:"league"`
											Sport      string `json:"sport"`
											GameID     int    `json:"gameId"`
											BetSide    string `json:"betSide"`
											BetType    string `json:"betType"`
											BetDetails string `json:"betDetails"`
										} `json:"tags"`
									} `json:"tracking"`
								} `json:"link"`
							} `json:"open"`
							Close struct {
								Line string `json:"line"`
								Odds string `json:"odds"`
								Link struct {
									Rel      []string `json:"rel"`
									Href     string   `json:"href"`
									Tracking struct {
										Campaign string `json:"campaign"`
										Tags     struct {
											League     string `json:"league"`
											Sport      string `json:"sport"`
											GameID     int    `json:"gameId"`
											BetSide    string `json:"betSide"`
											BetType    string `json:"betType"`
											BetDetails string `json:"betDetails"`
										} `json:"tags"`
									} `json:"tracking"`
								} `json:"link"`
							} `json:"close"`
						} `json:"over"`
						Under struct {
							Open struct {
								Line string `json:"line"`
								Odds string `json:"odds"`
								Link struct {
									Rel      []string `json:"rel"`
									Href     string   `json:"href"`
									Tracking struct {
										Campaign string `json:"campaign"`
										Tags     struct {
											League     string `json:"league"`
											Sport      string `json:"sport"`
											GameID     int    `json:"gameId"`
											BetSide    string `json:"betSide"`
											BetType    string `json:"betType"`
											BetDetails string `json:"betDetails"`
										} `json:"tags"`
									} `json:"tracking"`
								} `json:"link"`
							} `json:"open"`
							Close struct {
								Line string `json:"line"`
								Odds string `json:"odds"`
								Link struct {
									Rel      []string `json:"rel"`
									Href     string   `json:"href"`
									Tracking struct {
										Campaign string `json:"campaign"`
										Tags     struct {
											League     string `json:"league"`
											Sport      string `json:"sport"`
											GameID     int    `json:"gameId"`
											BetSide    string `json:"betSide"`
											BetType    string `json:"betType"`
											BetDetails string `json:"betDetails"`
										} `json:"tags"`
									} `json:"tracking"`
								} `json:"link"`
							} `json:"close"`
						} `json:"under"`
					} `json:"total"`
				} `json:"odds"`
				Situation struct {
					LastPlay struct {
						ID   string `json:"id"`
						Type struct {
							ID   string `json:"id"`
							Text string `json:"text"`
						} `json:"type"`
						Text      string `json:"text"`
						ShortText string `json:"shortText"`
						Period    struct {
							Number int `json:"number"`
						} `json:"period"`
						Clock struct {
							Value        int    `json:"value"`
							DisplayValue string `json:"displayValue"`
						} `json:"clock"`
						Team struct {
							ID string `json:"id"`
						} `json:"team"`
						ScoreValue       int `json:"scoreValue"`
						AthletesInvolved []struct {
							ID          string `json:"id"`
							DisplayName string `json:"displayName"`
							ShortName   string `json:"shortName"`
							FullName    string `json:"fullName"`
							Jersey      string `json:"jersey"`
							Headshot    string `json:"headshot"`
							Links       []struct {
								Rel      []string `json:"rel"`
								Href     string   `json:"href"`
								IsHidden bool     `json:"isHidden"`
							} `json:"links"`
							Position string `json:"position"`
						} `json:"athletesInvolved"`
						Probability struct {
							AwayWinPercentage float64 `json:"awayWinPercentage"`
							HomeWinPercentage float64 `json:"homeWinPercentage"`
							TiePercentage     float64 `json:"tiePercentage"`
						} `json:"probability"`
					} `json:"lastPlay"`
				} `json:"situation,omitempty"`
				FullStatus struct {
					Clock        float64 `json:"clock"`
					DisplayClock string  `json:"displayClock"`
					Period       int     `json:"period"`
					Type         struct {
						ID          string `json:"id"`
						Name        string `json:"name"`
						State       string `json:"state"`
						Completed   bool   `json:"completed"`
						Description string `json:"description"`
						Detail      string `json:"detail"`
						ShortDetail string `json:"shortDetail"`
					} `json:"type"`
				} `json:"fullStatus"`
				Competitors []struct {
					ID             string `json:"id"`
					UID            string `json:"uid"`
					Type           string `json:"type"`
					Order          int    `json:"order"`
					HomeAway       string `json:"homeAway"`
					Winner         bool   `json:"winner"`
					DisplayName    string `json:"displayName"`
					Name           string `json:"name"`
					Abbreviation   string `json:"abbreviation"`
					Location       string `json:"location"`
					Color          string `json:"color"`
					AlternateColor string `json:"alternateColor"`
					Score          string `json:"score"`
					Group          string `json:"group"`
					Record         string `json:"record"`
					Logo           string `json:"logo"`
					LogoDark       string `json:"logoDark"`
				} `json:"competitors"`
				AppLinks []struct {
					Rel       []string `json:"rel"`
					Href      string   `json:"href"`
					Text      string   `json:"text"`
					ShortText string   `json:"shortText"`
				} `json:"appLinks"`
				Video struct {
					Source      string `json:"source"`
					ID          int    `json:"id"`
					Headline    string `json:"headline"`
					Caption     string `json:"caption"`
					Description string `json:"description"`
					Premium     bool   `json:"premium"`
					Ad          struct {
						Sport  string `json:"sport"`
						Bundle string `json:"bundle"`
					} `json:"ad"`
					Tracking struct {
						SportName    string `json:"sportName"`
						LeagueName   string `json:"leagueName"`
						CoverageType string `json:"coverageType"`
						TrackingName string `json:"trackingName"`
						TrackingID   string `json:"trackingId"`
					} `json:"tracking"`
					CerebroID           string    `json:"cerebroId"`
					PccID               string    `json:"pccId"`
					LastModified        time.Time `json:"lastModified"`
					OriginalPublishDate time.Time `json:"originalPublishDate"`
					TimeRestrictions    struct {
						EmbargoDate    time.Time `json:"embargoDate"`
						ExpirationDate time.Time `json:"expirationDate"`
					} `json:"timeRestrictions"`
					DeviceRestrictions struct {
						Type    string   `json:"type"`
						Devices []string `json:"devices"`
					} `json:"deviceRestrictions"`
					GeoRestrictions struct {
						Type      string   `json:"type"`
						Countries []string `json:"countries"`
					} `json:"geoRestrictions"`
					Syndicatable bool `json:"syndicatable"`
					Duration     int  `json:"duration"`
					Categories   []struct {
						ID          int    `json:"id"`
						Description string `json:"description"`
						Type        string `json:"type"`
						SportID     int    `json:"sportId"`
						LeagueID    int    `json:"leagueId,omitempty"`
						League      struct {
							ID          int    `json:"id"`
							Description string `json:"description"`
							Links       struct {
								API struct {
									Leagues struct {
										Href string `json:"href"`
									} `json:"leagues"`
								} `json:"api"`
								Web struct {
									Leagues struct {
										Href string `json:"href"`
									} `json:"leagues"`
								} `json:"web"`
								Mobile struct {
									Leagues struct {
										Href string `json:"href"`
									} `json:"leagues"`
								} `json:"mobile"`
							} `json:"links"`
						} `json:"league,omitempty"`
						UID     string `json:"uid,omitempty"`
						TopicID int    `json:"topicId,omitempty"`
						TeamID  int    `json:"teamId,omitempty"`
						Team    struct {
							ID          int    `json:"id"`
							Description string `json:"description"`
							Links       struct {
								API struct {
									Teams struct {
										Href string `json:"href"`
									} `json:"teams"`
								} `json:"api"`
								Web struct {
									Teams struct {
										Href string `json:"href"`
									} `json:"teams"`
								} `json:"web"`
								Mobile struct {
									Teams struct {
										Href string `json:"href"`
									} `json:"teams"`
								} `json:"mobile"`
							} `json:"links"`
						} `json:"team,omitempty"`
						AthleteID int `json:"athleteId,omitempty"`
						Athlete   struct {
							ID          int    `json:"id"`
							Description string `json:"description"`
							Links       struct {
								API struct {
									Athletes struct {
										Href string `json:"href"`
									} `json:"athletes"`
								} `json:"api"`
								Web struct {
									Athletes struct {
										Href string `json:"href"`
									} `json:"athletes"`
								} `json:"web"`
								Mobile struct {
									Athletes struct {
										Href string `json:"href"`
									} `json:"athletes"`
								} `json:"mobile"`
							} `json:"links"`
						} `json:"athlete,omitempty"`
					} `json:"categories"`
					GameID       int   `json:"gameId"`
					Keywords     []any `json:"keywords"`
					PosterImages struct {
						Default struct {
							Href   string `json:"href"`
							Width  int    `json:"width"`
							Height int    `json:"height"`
						} `json:"default"`
						Full struct {
							Href string `json:"href"`
						} `json:"full"`
						Wide struct {
							Href string `json:"href"`
						} `json:"wide"`
						Square struct {
							Href string `json:"href"`
						} `json:"square"`
					} `json:"posterImages"`
					Images []struct {
						Name    string `json:"name"`
						URL     string `json:"url"`
						Alt     string `json:"alt"`
						Caption string `json:"caption"`
						Credit  string `json:"credit"`
						Width   int    `json:"width"`
						Height  int    `json:"height"`
					} `json:"images"`
					Thumbnail string `json:"thumbnail"`
					Links     struct {
						API struct {
							Self struct {
								Href string `json:"href"`
							} `json:"self"`
							Artwork struct {
								Href string `json:"href"`
							} `json:"artwork"`
						} `json:"api"`
						Web struct {
							Href  string `json:"href"`
							Short struct {
								Href string `json:"href"`
							} `json:"short"`
							Self struct {
								Href string `json:"href"`
							} `json:"self"`
						} `json:"web"`
						Source struct {
							Mezzanine struct {
								Href string `json:"href"`
							} `json:"mezzanine"`
							Flash struct {
								Href string `json:"href"`
							} `json:"flash"`
							Hds struct {
								Href string `json:"href"`
							} `json:"hds"`
							Hls struct {
								Href string `json:"href"`
								Hd   struct {
									Href string `json:"href"`
								} `json:"HD"`
							} `json:"HLS"`
							Hd struct {
								Href string `json:"href"`
							} `json:"HD"`
							Full struct {
								Href string `json:"href"`
							} `json:"full"`
							Href string `json:"href"`
						} `json:"source"`
						Mobile struct {
							Alert struct {
								Href string `json:"href"`
							} `json:"alert"`
							Source struct {
								Href string `json:"href"`
							} `json:"source"`
							Href      string `json:"href"`
							Streaming struct {
								Href string `json:"href"`
							} `json:"streaming"`
							ProgressiveDownload struct {
								Href string `json:"href"`
							} `json:"progressiveDownload"`
						} `json:"mobile"`
					} `json:"links"`
				} `json:"video,omitempty"`
			} `json:"events"`
		} `json:"leagues"`
	} `json:"sports"`
	Zipcodes struct {
	} `json:"zipcodes"`
}

// GameResult 用于存储和排序比赛信息// GameInfo 用于存储和排序比赛信息
type GameInfo struct {
	Description string
	GameTime    time.Time
}
