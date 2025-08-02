package rl

// GetRealisticTrainingData returns a comprehensive set of realistic training examples
func GetRealisticTrainingData() []TrainingExample {
	return []TrainingExample{
		// Technical Documentation Examples
		{
			ID: "tech_doc_1",
			Text: `The Redis persistence mechanism offers two distinct approaches: RDB (Redis Database) snapshots 
and AOF (Append Only File) logging. RDB performs point-in-time snapshots of your dataset at specified 
intervals, while AOF logs every write operation received by the server. These methods can be used 
independently or combined for maximum data safety. The trade-off involves balancing performance 
impact against data durability requirements.`,
			TaskType: "technical_analysis",
			Expected: map[string]interface{}{
				"entities": []string{"Redis", "RDB", "AOF", "Redis Database", "Append Only File"},
				"has_code": false,
				"readability_score": 65.0,
				"keywords": []string{"persistence", "snapshots", "logging", "performance", "durability"},
				"complexity": "medium",
			},
			Difficulty: 0.6,
		},
		
		// Code Sample with Documentation
		{
			ID: "code_sample_1",
			Text: `def calculate_fibonacci(n):
    """
    Calculate the nth Fibonacci number using dynamic programming.
    Time complexity: O(n), Space complexity: O(1)
    """
    if n <= 1:
        return n
    
    prev, curr = 0, 1
    for i in range(2, n + 1):
        prev, curr = curr, prev + curr
    
    return curr

# Example usage:
# print(calculate_fibonacci(10))  # Output: 55`,
			TaskType: "code_analysis",
			Expected: map[string]interface{}{
				"has_code": true,
				"language": "python",
				"entities": []string{"fibonacci", "dynamic programming"},
				"code_complexity": "low",
				"has_comments": true,
				"function_count": 1,
			},
			Difficulty: 0.7,
		},
		
		// Research Paper Abstract
		{
			ID: "research_abstract_1",
			Text: `Abstract: We present a novel approach to multi-task learning in neural networks that 
significantly improves performance on diverse NLP tasks. Our method, termed Adaptive Task 
Prioritization (ATP), dynamically adjusts task weights during training based on gradient 
similarity and task-specific loss trajectories. Experiments on the GLUE benchmark demonstrate 
that ATP achieves state-of-the-art results on 7 out of 9 tasks, with an average improvement 
of 3.2% over baseline MTL approaches. Furthermore, our analysis reveals that ATP effectively 
mitigates negative transfer between dissimilar tasks while promoting positive transfer among 
related tasks.`,
			TaskType: "academic_analysis",
			Expected: map[string]interface{}{
				"entities": []string{"neural networks", "NLP", "Adaptive Task Prioritization", "ATP", "GLUE"},
				"sentiment": "positive",
				"keywords": []string{"multi-task learning", "gradient", "benchmark", "state-of-the-art"},
				"readability_score": 45.0,
				"has_statistics": true,
			},
			Difficulty: 0.8,
		},
		
		// Business Email
		{
			ID: "business_email_1",
			Text: `Subject: Q3 Product Roadmap Review - Action Required

Hi Team,

Following our strategic planning session last week, I'm sharing the updated Q3 roadmap priorities:

1. Customer Portal v2.0 - Target: July 15th
   - Enhanced dashboard with real-time analytics
   - Mobile-responsive design
   - SSO integration

2. API Rate Limiting - Target: August 1st
   - Implement tiered usage plans
   - Add monitoring and alerting

3. Performance Optimization - Ongoing
   - Database query optimization (30% improvement target)
   - CDN implementation for static assets

Please review and provide feedback by EOD Friday. We'll finalize during Monday's standup.

Best regards,
Sarah Chen
Product Manager`,
			TaskType: "business_communication",
			Expected: map[string]interface{}{
				"entities": []string{"Customer Portal", "API", "SSO", "CDN", "Sarah Chen"},
				"sentiment": "neutral",
				"has_action_items": true,
				"readability_score": 75.0,
				"format_type": "email",
				"has_deadline": true,
			},
			Difficulty: 0.5,
		},
		
		// News Article Excerpt
		{
			ID: "news_article_1",
			Text: `The European Central Bank announced today a 0.25 percentage point increase in interest rates, 
marking the tenth consecutive hike in its aggressive campaign against inflation. The deposit 
rate now stands at 4%, the highest level since the 2008 financial crisis. ECB President 
Christine Lagarde emphasized that "inflation remains too high for too long," citing persistent 
core inflation at 5.3%. Financial markets responded positively, with the euro gaining 0.8% 
against the dollar. However, concerns mount about the impact on mortgage holders and business 
investment across the eurozone's struggling economies.`,
			TaskType: "news_analysis",
			Expected: map[string]interface{}{
				"entities": []string{"European Central Bank", "ECB", "Christine Lagarde"},
				"sentiment": "mixed",
				"keywords": []string{"interest rates", "inflation", "financial crisis", "euro", "mortgage"},
				"has_quotes": true,
				"has_statistics": true,
				"readability_score": 60.0,
			},
			Difficulty: 0.6,
		},
		
		// Social Media Post
		{
			ID: "social_media_1",
			Text: `🚀 Excited to announce that our team just shipped the biggest update to @CloudSyncApp yet! 

✨ What's new:
- 10x faster file uploads 
- End-to-end encryption 
- Collaborative folders
- Dark mode (finally!)

Thank you to our amazing beta testers who provided invaluable feedback. Special shoutout to 
@DevSarah and @TechMike for finding those edge cases!

Try it out and let us know what you think 👇

#CloudStorage #ProductLaunch #TechNews #StartupLife`,
			TaskType: "social_media_analysis",
			Expected: map[string]interface{}{
				"entities": []string{"CloudSyncApp", "DevSarah", "TechMike"},
				"sentiment": "positive",
				"has_emoji": true,
				"has_hashtags": true,
				"engagement_potential": "high",
				"readability_score": 85.0,
			},
			Difficulty: 0.4,
		},
		
		// Legal Document Excerpt
		{
			ID: "legal_doc_1",
			Text: `WHEREAS, the Party of the First Part (hereinafter "Licensor") owns certain intellectual 
property rights in and to the software known as "DataFlow Analytics Suite" (the "Software"); 
and WHEREAS, the Party of the Second Part (hereinafter "Licensee") desires to obtain a 
non-exclusive, non-transferable license to use the Software subject to the terms and 
conditions set forth herein; NOW, THEREFORE, in consideration of the mutual covenants and 
agreements contained herein, and for other good and valuable consideration, the receipt 
and sufficiency of which are hereby acknowledged, the parties agree as follows:`,
			TaskType: "legal_analysis",
			Expected: map[string]interface{}{
				"entities": []string{"DataFlow Analytics Suite", "Licensor", "Licensee"},
				"document_type": "license_agreement",
				"readability_score": 25.0,
				"formal_language": true,
				"has_legal_terms": true,
			},
			Difficulty: 0.9,
		},
		
		// Medical/Healthcare Text
		{
			ID: "medical_1",
			Text: `Patient presented with acute onset chest pain, radiating to left arm, accompanied by 
diaphoresis and dyspnea. ECG revealed ST-segment elevation in leads II, III, and aVF, 
consistent with inferior wall MI. Troponin I elevated at 2.5 ng/mL (normal <0.04). 
Initiated standard ACS protocol: aspirin 325mg, clopidogrel 600mg, and heparin bolus. 
Urgent cardiac catheterization revealed 95% occlusion of RCA. Successfully performed PCI 
with drug-eluting stent placement. Post-procedure TIMI 3 flow achieved.`,
			TaskType: "medical_analysis",
			Expected: map[string]interface{}{
				"entities": []string{"ECG", "ST-segment", "MI", "Troponin I", "ACS", "RCA", "PCI", "TIMI"},
				"medical_terms": true,
				"urgency_level": "high",
				"has_measurements": true,
				"readability_score": 20.0,
			},
			Difficulty: 0.9,
		},
		
		// Recipe/Instructional Text
		{
			ID: "recipe_1",
			Text: `Classic French Onion Soup

Ingredients:
- 6 large yellow onions, thinly sliced
- 4 tbsp butter
- 2 cloves garlic, minced
- 1/2 cup dry sherry
- 8 cups beef broth
- 2 bay leaves
- Fresh thyme
- Gruyere cheese, grated
- Baguette slices

Instructions:
1. Caramelize onions in butter over medium-low heat for 45 minutes, stirring occasionally
2. Add garlic and cook for 2 minutes
3. Deglaze with sherry, scraping up browned bits
4. Add broth, bay leaves, and thyme; simmer for 20 minutes
5. Top with toasted baguette slices and cheese, broil until golden

Prep time: 15 minutes | Cook time: 1 hour 10 minutes | Serves: 6`,
			TaskType: "instructional_analysis",
			Expected: map[string]interface{}{
				"entities": []string{"French Onion Soup", "Gruyere cheese"},
				"has_ingredients_list": true,
				"has_instructions": true,
				"has_timing_info": true,
				"readability_score": 80.0,
				"format_type": "recipe",
			},
			Difficulty: 0.3,
		},
		
		// Error Log/Technical Output
		{
			ID: "error_log_1",
			Text: `[2024-03-14 15:32:41.823] ERROR: Database connection failed
java.sql.SQLException: Connection refused: connect
    at com.mysql.jdbc.StandardConnection.connect(StandardConnection.java:203)
    at com.example.dao.UserDAO.getConnection(UserDAO.java:45)
    at com.example.service.AuthService.authenticate(AuthService.java:78)
Caused by: java.net.ConnectException: Connection timed out
    ... 12 more
[2024-03-14 15:32:41.825] WARN: Falling back to cache for user authentication
[2024-03-14 15:32:41.830] INFO: Cache hit for user_id: 12345
[2024-03-14 15:32:41.832] INFO: Authentication successful via cache`,
			TaskType: "log_analysis",
			Expected: map[string]interface{}{
				"has_error": true,
				"error_type": "SQLException",
				"has_timestamps": true,
				"log_levels": []string{"ERROR", "WARN", "INFO"},
				"has_stack_trace": true,
				"readability_score": 30.0,
			},
			Difficulty: 0.7,
		},
		
		// Marketing Copy
		{
			ID: "marketing_1",
			Text: `Transform Your Morning Routine with BrewMaster Pro

Wake up to barista-quality coffee without leaving your home. The BrewMaster Pro combines 
cutting-edge extraction technology with intuitive one-touch controls, delivering the perfect 
cup every time.

✓ Precision temperature control (195-205°F)
✓ 15-bar pressure system
✓ Built-in grinder with 30 settings
✓ Self-cleaning function

Limited time offer: Save $100 and get free shipping. Plus, receive a complimentary bag of 
our signature Ethiopian blend with your purchase.

Join 50,000+ coffee enthusiasts who've made the switch. Order now and taste the difference!`,
			TaskType: "marketing_analysis",
			Expected: map[string]interface{}{
				"entities": []string{"BrewMaster Pro", "Ethiopian blend"},
				"sentiment": "positive",
				"has_call_to_action": true,
				"has_features_list": true,
				"persuasion_techniques": []string{"social_proof", "urgency", "benefits"},
				"readability_score": 70.0,
			},
			Difficulty: 0.5,
		},
		
		// Scientific/Environmental Report
		{
			ID: "scientific_1",
			Text: `The Arctic permafrost contains approximately 1,700 billion metric tons of carbon, nearly 
twice the amount currently in the atmosphere. Recent studies indicate that permafrost thaw 
rates have accelerated by 240% since the 1990s, with mean annual ground temperatures 
increasing by 2.3°C. This positive feedback loop releases both CO2 and methane, with 
methane emissions being particularly concerning due to its 28x greater warming potential 
over a 100-year period. Model projections suggest that 40-60% of permafrost area could 
thaw by 2100 under RCP 8.5 scenarios, potentially releasing 130-160 Gt of carbon.`,
			TaskType: "scientific_analysis",
			Expected: map[string]interface{}{
				"entities": []string{"Arctic", "CO2", "methane", "RCP 8.5"},
				"has_statistics": true,
				"has_measurements": true,
				"scientific_terms": true,
				"sentiment": "negative",
				"readability_score": 40.0,
			},
			Difficulty: 0.8,
		},
	}
}