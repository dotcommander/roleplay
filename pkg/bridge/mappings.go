package bridge

// TraitMapping defines how a personality trait maps to OCEAN dimensions.
// Values range from -0.5 to +0.5, where:
// - Positive values increase the dimension
// - Negative values decrease the dimension
// - Zero means no effect on that dimension
type TraitMapping struct {
	Openness          float64
	Conscientiousness float64
	Extraversion      float64
	Agreeableness     float64
	Neuroticism       float64
}

// TraitMappings contains all trait to OCEAN mappings.
type TraitMappings struct {
	Traits map[string]TraitMapping
}

// GetDefaultMappings returns the default trait mappings database.
func GetDefaultMappings() *TraitMappings {
	return &TraitMappings{
		Traits: map[string]TraitMapping{
			// Openness traits
			"creative":     {Openness: 0.4},
			"imaginative":  {Openness: 0.4},
			"artistic":     {Openness: 0.3},
			"curious":      {Openness: 0.3},
			"adventurous":  {Openness: 0.3, Extraversion: 0.2},
			"open-minded":  {Openness: 0.4},
			"intellectual": {Openness: 0.3},
			"innovative":   {Openness: 0.4},
			"philosophical": {Openness: 0.3},
			"unconventional": {Openness: 0.3},
			"traditional":  {Openness: -0.3},
			"conservative": {Openness: -0.3},
			"practical":    {Openness: -0.2, Conscientiousness: 0.2},
			"routine":      {Openness: -0.2},

			// Conscientiousness traits
			"organized":     {Conscientiousness: 0.4},
			"disciplined":   {Conscientiousness: 0.4},
			"responsible":   {Conscientiousness: 0.3},
			"reliable":      {Conscientiousness: 0.3},
			"hardworking":   {Conscientiousness: 0.3},
			"meticulous":    {Conscientiousness: 0.4},
			"punctual":      {Conscientiousness: 0.3},
			"perfectionist": {Conscientiousness: 0.3, Neuroticism: 0.2},
			"methodical":    {Conscientiousness: 0.3},
			"careful":       {Conscientiousness: 0.3},
			"impulsive":     {Conscientiousness: -0.3, Neuroticism: 0.1},
			"careless":      {Conscientiousness: -0.3},
			"disorganized":  {Conscientiousness: -0.4},
			"lazy":          {Conscientiousness: -0.3},
			"procrastinator": {Conscientiousness: -0.3},

			// Extraversion traits
			"outgoing":      {Extraversion: 0.4},
			"sociable":      {Extraversion: 0.4},
			"talkative":     {Extraversion: 0.3},
			"energetic":     {Extraversion: 0.3},
			"enthusiastic":  {Extraversion: 0.3},
			"assertive":     {Extraversion: 0.3},
			"charismatic":   {Extraversion: 0.3, Agreeableness: 0.1},
			"gregarious":    {Extraversion: 0.4},
			"extraverted":   {Extraversion: 0.4},
			"extroverted":   {Extraversion: 0.4},
			"shy":           {Extraversion: -0.3, Neuroticism: 0.1},
			"introverted":   {Extraversion: -0.4},
			"reserved":      {Extraversion: -0.3},
			"quiet":         {Extraversion: -0.3},
			"solitary":      {Extraversion: -0.3},
			"withdrawn":     {Extraversion: -0.3, Neuroticism: 0.1},

			// Agreeableness traits
			"friendly":      {Agreeableness: 0.4},
			"kind":          {Agreeableness: 0.4},
			"compassionate": {Agreeableness: 0.4},
			"helpful":       {Agreeableness: 0.3},
			"cooperative":   {Agreeableness: 0.3},
			"trusting":      {Agreeableness: 0.3},
			"empathetic":    {Agreeableness: 0.4},
			"generous":      {Agreeableness: 0.3},
			"warm":          {Agreeableness: 0.3},
			"caring":        {Agreeableness: 0.4},
			"agreeable":     {Agreeableness: 0.3},
			"cynical":       {Agreeableness: -0.3},
			"critical":      {Agreeableness: -0.2},
			"competitive":   {Agreeableness: -0.2, Conscientiousness: 0.1},
			"aggressive":    {Agreeableness: -0.3, Neuroticism: 0.1},
			"argumentative": {Agreeableness: -0.3},
			"selfish":       {Agreeableness: -0.4},
			"manipulative":  {Agreeableness: -0.4},

			// Neuroticism traits
			"anxious":       {Neuroticism: 0.4},
			"nervous":       {Neuroticism: 0.3},
			"stressed":      {Neuroticism: 0.3},
			"emotional":     {Neuroticism: 0.3},
			"moody":         {Neuroticism: 0.3},
			"temperamental": {Neuroticism: 0.3},
			"sensitive":     {Neuroticism: 0.2, Agreeableness: 0.1},
			"worrier":       {Neuroticism: 0.4},
			"insecure":      {Neuroticism: 0.3},
			"neurotic":      {Neuroticism: 0.4},
			"calm":          {Neuroticism: -0.3},
			"stable":        {Neuroticism: -0.3},
			"confident":     {Neuroticism: -0.3, Extraversion: 0.1},
			"relaxed":       {Neuroticism: -0.3},
			"resilient":     {Neuroticism: -0.3, Conscientiousness: 0.1},
			"secure":        {Neuroticism: -0.3},

			// Complex/Combined traits
			"analytical":    {Conscientiousness: 0.2, Openness: 0.2},
			"ambitious":     {Conscientiousness: 0.3, Extraversion: 0.2},
			"independent":   {Extraversion: -0.1, Agreeableness: -0.1},
			"leader":        {Extraversion: 0.3, Conscientiousness: 0.2},
			"rebellious":    {Agreeableness: -0.2, Openness: 0.2},
			"strategic":     {Conscientiousness: 0.3, Openness: 0.1},
			"spontaneous":   {Openness: 0.3, Conscientiousness: -0.2},
			"eccentric":     {Openness: 0.3, Agreeableness: -0.1},
			"charming":      {Extraversion: 0.2, Agreeableness: 0.2},
			"witty":         {Openness: 0.2, Extraversion: 0.2},
			"sarcastic":     {Agreeableness: -0.2, Openness: 0.1},
			"loyal":         {Agreeableness: 0.3, Conscientiousness: 0.2},
			"honest":        {Agreeableness: 0.2, Conscientiousness: 0.2},
			"patient":       {Agreeableness: 0.2, Neuroticism: -0.2},
			"optimistic":    {Neuroticism: -0.2, Extraversion: 0.2},
			"pessimistic":   {Neuroticism: 0.2, Extraversion: -0.1},
			"logical":       {Conscientiousness: 0.2, Agreeableness: -0.1},
			"intuitive":     {Openness: 0.3},
			"determined":    {Conscientiousness: 0.3},
			"flexible":      {Openness: 0.2, Agreeableness: 0.1},
			"stubborn":      {Agreeableness: -0.2, Conscientiousness: 0.1},
		},
	}
}

// AddCustomMapping adds a custom trait mapping.
func (tm *TraitMappings) AddCustomMapping(trait string, mapping TraitMapping) {
	if tm.Traits == nil {
		tm.Traits = make(map[string]TraitMapping)
	}
	tm.Traits[trait] = mapping
}

// GetMapping retrieves a trait mapping.
func (tm *TraitMappings) GetMapping(trait string) (TraitMapping, bool) {
	mapping, exists := tm.Traits[trait]
	return mapping, exists
}

// MergeWith merges another set of mappings into this one.
func (tm *TraitMappings) MergeWith(other *TraitMappings) {
	if tm.Traits == nil {
		tm.Traits = make(map[string]TraitMapping)
	}
	for trait, mapping := range other.Traits {
		tm.Traits[trait] = mapping
	}
}