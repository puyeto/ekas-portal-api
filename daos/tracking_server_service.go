package daos

// TrackingServerServiceDAO persists arackingServerService data in database
type TrackingServerServiceDAO struct{}

// NewTrackingServerServiceDAO creates a new TrackingServerServiceDAO
func NewTrackingServerServiceDAO() *TrackingServerServiceDAO {
	return &TrackingServerServiceDAO{}
}
