package voters

import (
	"errors"
	"time"
	"voter-api/voters/utils"
)

type voterPoll struct {
	PollID   uint      `json:"pollId"`
	VoteDate time.Time `json:"voteDate"`
}

type VoteHistory []voterPoll

type Voter struct {
	VoterID     uint        `json:"voterId"`
	FirstName   string      `json:"firstName"`
	LastName    string      `json:"lastName"`
	VoteHistory VoteHistory `json:"votes"`
}

type VoterList struct {
	Voters map[uint]Voter
}

func NewVoterList() (*VoterList, error) {
	list := &VoterList{
		Voters: make(map[uint]Voter),
	}

	return list, nil
}

func NewVoter(id uint, firstName string, lastName string) (Voter, error) {
	voter := Voter{
		VoterID:     id,
		FirstName:   firstName,
		LastName:    lastName,
		VoteHistory: make([]voterPoll, 0),
	}

	return voter, nil
}

func (v *VoterList) AddVoter(id uint, voter Voter) (Voter, error) {
	_, err := v.GetVoter(id)

	if err == nil {
		return Voter{}, errors.New("Voter already exists")
	}

	voter.VoterID = id

	v.Voters[id] = voter

	return v.Voters[id], nil
}

func (v *VoterList) GetVoter(id uint) (*Voter, error) {
	voter, ok := v.Voters[id]

	if !ok {
		return &Voter{}, errors.New("Voter not found")
	}

	return &voter, nil
}

func (v *VoterList) UpdateVoter(id uint, voter Voter) (Voter, error) {
	existingVoter, err := v.GetVoter(id)

	if err != nil {
		return Voter{}, err
	}

	existingVoter.FirstName = voter.FirstName
	existingVoter.LastName = voter.LastName

	v.Voters[id] = *existingVoter

	return *existingVoter, nil
}

func (v *VoterList) DeleteVoter(id uint) error {
	_, err := v.GetVoter(id)

	if err != nil {
		return err
	}

	delete(v.Voters, id)

	return nil
}

func (v *VoterList) GetVoterHistory(id uint) ([]voterPoll, error) {
	voter, err := v.GetVoter(id)

	if err != nil {
		return nil, err
	}

	return voter.VoteHistory, nil
}

func (v *VoterList) GetVoterVote(id uint, pollId uint) (*voterPoll, int, error) {
	voter, err := v.GetVoter(id)

	if err != nil {
		return &voterPoll{}, 0, err
	}

	for i, vote := range voter.VoteHistory {
		if vote.PollID == pollId {
			return &vote, i, nil
		}
	}

	return &voterPoll{}, 0, errors.New("Vote not found")
}

func (v *VoterList) CreateVoterVote(id uint, pollId uint) (voterPoll, error) {
	voter, voterErr := v.GetVoter(id)

	if voterErr != nil {
		return voterPoll{}, voterErr
	}

	_, _, voteErr := v.GetVoterVote(id, pollId)

	if voteErr == nil {
		return voterPoll{}, errors.New("Vote for that poll already exists")
	}

	newVoterPoll := voterPoll{
		PollID:   pollId,
		VoteDate: time.Now(),
	}

	voter.VoteHistory = append(voter.VoteHistory, newVoterPoll)

	v.Voters[id] = *voter

	return newVoterPoll, nil
}

func (v *VoterList) UpdateVoterVote(id uint, pollId uint) (voterPoll, error) {
	voter, voterErr := v.GetVoter(id)

	if voterErr != nil {
		return voterPoll{}, voterErr
	}

	vote, voteIndex, voteErr := v.GetVoterVote(id, pollId)

	if voteErr != nil {
		return voterPoll{}, voteErr
	}

	vote.VoteDate = time.Now()

	voter.VoteHistory[voteIndex] = *vote

	voter.VoteHistory = utils.ShiftEnd(voter.VoteHistory, voteIndex)

	v.Voters[id] = *voter

	return *vote, nil
}

func (v *VoterList) DeleteVoterVote(id uint, pollId uint) error {
	voter, voterErr := v.GetVoter(id)

	if voterErr != nil {
		return voterErr
	}

	_, voteIndex, voteErr := v.GetVoterVote(id, pollId)

	if voteErr != nil {
		return voteErr
	}

	voter.VoteHistory = append(voter.VoteHistory[:voteIndex], voter.VoteHistory[voteIndex+1:]...)

	v.Voters[id] = *voter

	return nil
}

func (v *VoterList) GetVoters() ([]Voter, error) {
	voters := make([]Voter, 0)

	for _, voter := range v.Voters {
		voters = append(voters, voter)
	}

	return voters, nil
}
