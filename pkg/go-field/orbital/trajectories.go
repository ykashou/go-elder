package orbital

type TrajectoryCalculator struct {
	InitialConditions OrbitalMechanics
	TimeStep          float64
}

func (tc *TrajectoryCalculator) ComputeTrajectory(duration float64) []Vector3D {
	trajectory := []Vector3D{}
	current := tc.InitialConditions
	
	for t := 0.0; t < duration; t += tc.TimeStep {
		trajectory = append(trajectory, current.Position)
		current.Position.X += current.Velocity.X * tc.TimeStep
		current.Position.Y += current.Velocity.Y * tc.TimeStep
		current.Position.Z += current.Velocity.Z * tc.TimeStep
	}
	
	return trajectory
}
