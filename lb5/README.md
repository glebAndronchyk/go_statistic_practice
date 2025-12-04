# Lab Work 7: Pearson's Chi-Squared Goodness-of-Fit Test

## Table of Contents
1. [Overview](#overview)
2. [Theoretical Background](#theoretical-background)
3. [Implementation Details](#implementation-details)
4. [Code Structure](#code-structure)
5. [Mathematical Formulas](#mathematical-formulas)
6. [Usage Examples](#usage-examples)
7. [Interpretation of Results](#interpretation-of-results)

---

## Overview

This program implements **Pearson's Chi-Squared (χ²) Goodness-of-Fit Test** to determine whether empirical data follows a specific theoretical distribution. The program tests two distributions:
- **Normal (Gaussian) Distribution**
- **Uniform Distribution**

### Purpose
Given a sample of data, we want to answer the question: "Does this data come from a normal/uniform distribution?"

The test uses statistical hypothesis testing:
- **H₀ (Null Hypothesis)**: The data follows the specified distribution
- **H₁ (Alternative Hypothesis)**: The data does NOT follow the specified distribution

---

## Theoretical Background

### What is the Chi-Squared Test?

The Chi-Squared test compares **observed frequencies** (actual data) with **expected frequencies** (theoretical distribution) to determine if the differences are due to random chance or indicate the data doesn't fit the distribution.

### The Chi-Squared Statistic

The test statistic is calculated as:

```
χ² = Σ [(nᵢ - nᵢᵀ)² / nᵢᵀ]
```

Where:
- `nᵢ` = observed frequency in bin i
- `nᵢᵀ` = theoretical (expected) frequency in bin i
- The sum is over all bins

**Interpretation:**
- **Small χ²**: Observed and expected frequencies are similar → data fits the distribution
- **Large χ²**: Significant differences → data doesn't fit the distribution

### Decision Rule

1. Calculate empirical χ² from your data (χ²ₑₘₚ)
2. Find critical value χ²ₖᵣ for significance level α
3. Compare:
   - If **χ²ₑₘₚ < χ²ₖᵣ**: Accept H₀ (data fits distribution)
   - If **χ²ₑₘₚ ≥ χ²ₖᵣ**: Reject H₀ (data doesn't fit)

### Degrees of Freedom

The degrees of freedom (r) depends on:
```
r = m - k - 1
```

Where:
- `m` = number of bins (intervals)
- `k` = number of parameters estimated from the data
  - Normal distribution: k = 2 (mean and standard deviation)
  - Uniform distribution: k = 0 (no parameters estimated)

---

## Implementation Details

### 1. Data Input and Validation (`readInput`)

```go
func readInput() ([]int, float64)
```

**What it does:**
- Reads space-separated integers from console
- Validates each number is in range [0, 20]
- Skips invalid values with warnings
- Reads significance level α
- Validates α is between 0 and 1

**Why this matters:**
The Chi-Squared test requires:
- Sufficient sample size (typically n ≥ 30)
- Frequency in each bin ≥ 5
- Valid significance level (usually 0.01, 0.05, or 0.10)

---

### 2. Histogram Creation (`createHistogram`, `mergeBins`)

```go
func createHistogram(sample []int, min, max int) *Histogram
```

**What it does:**
1. Counts frequency of each value
2. Creates bins for each unique value
3. **Merges adjacent bins** if any has count < 5

**Why merging is crucial:**
The Chi-Squared test is **not reliable** when expected frequencies are too small. The rule of thumb:
- Each bin should have **at least 5 observations**
- If not, merge with adjacent bins

**Example:**
```
Before merging:
[0]: 2 observations  ← too small!
[1]: 3 observations  ← too small!
[2]: 8 observations

After merging:
[0-1]: 5 observations  ✓
[2]: 8 observations    ✓
```

---

### 3. Goldstein Approximation (`calculateChiSquaredCritical`)

```go
func calculateChiSquaredCritical(alpha float64, r int) float64
```

**What it does:**
Calculates the critical value χ²ₖᵣ without using lookup tables.

**The Formula:**
```
χ²α,n = n · [Σᵢ₌₀⁶ n^(-i/2) · dⁱ · (aᵢ + bᵢ/n + cᵢ/n²)]³
```

Where `d` is calculated based on α:
- For 0.5 ≤ α ≤ 0.999:
  ```
  d = 2.0637 · (ln(1/(1-α)) - 0.16)^0.4274 - 1.5774
  ```
- For 0.001 ≤ α < 0.5:
  ```
  d = -2.0637 · (ln(1/α) - 0.16)^0.4274 + 1.5774
  ```

**Coefficients:**
The coefficients (a₀-a₆, b₀-b₆, c₀-c₆) are from the table in the PDF and stored in `goldsteinCoeffs`.

**Why use approximation?**
- Tables don't cover all possible α and r values
- Allows testing at any significance level
- More flexible than table lookup

---

### 4. Normal Distribution Test (`testNormalDistribution`)

```go
func testNormalDistribution(hist *Histogram, alpha float64, n int)
```

**Step-by-step process:**

#### Step 1: Estimate Parameters
Calculate sample mean and standard deviation:

```
x̄ = (1/n) · Σxᵢ

s = √[(1/(n-1)) · Σ(xᵢ - x̄)²]
```

**Why estimate?**
We don't know the "true" parameters of the population, so we estimate from the sample.

#### Step 2: Calculate Theoretical Probabilities
For each bin [a, b], calculate the probability using the normal CDF:

```
P(a < x < b) = Φ((b - μ)/σ) - Φ((a - μ)/σ)
```

Where Φ is the cumulative distribution function of the standard normal distribution.

**Implementation:**
```go
prob := normalCDF((upper-mean)/stdDev) - normalCDF((lower-mean)/stdDev)
```

The `normalCDF` function uses the error function:
```go
func normalCDF(x float64) float64 {
    return 0.5 * (1.0 + math.Erf(x/math.Sqrt2))
}
```

#### Step 3: Calculate Expected Frequencies
```
nᵢᵀ = n · P(bin i)
```

#### Step 4: Compute χ² Statistic
```go
chiSquared := 0.0
for each bin {
    observed := bin.Count
    expected := theoreticalFreq[i]
    contribution := (observed - expected)² / expected
    chiSquared += contribution
}
```

#### Step 5: Determine Degrees of Freedom
```
r = m - k - 1 = m - 2 - 1 = m - 3
```

Why k=2? We estimated 2 parameters (mean and std dev).

#### Step 6: Compare and Decide
```go
if chiSquared < chiSquaredCritical {
    // Accept H₀: data follows normal distribution
} else {
    // Reject H₀: data doesn't follow normal distribution
}
```

---

### 5. Uniform Distribution Test (`testUniformDistribution`)

```go
func testUniformDistribution(hist *Histogram, alpha float64, n int)
```

**Step-by-step process:**

#### Step 1: Determine Range
For uniform distribution U[a, b], we use:
- a = minimum value in data
- b = maximum value in data

**Example:** If data ranges from 9 to 12, we test U[9, 12].

#### Step 2: Calculate Theoretical Probabilities
For a bin covering values [lower, upper]:

```
Bin size = upper - lower + 1
Range size = max - min + 1
Probability = Bin size / Range size
```

**Example:**
For U[9, 12] (range size = 4):
- Bin [9]: probability = 1/4 = 0.25
- Bin [10]: probability = 1/4 = 0.25
- Bin [11]: probability = 1/4 = 0.25
- Bin [12]: probability = 1/4 = 0.25

#### Step 3: Calculate Expected Frequencies
```
nᵢᵀ = n · P(bin i)
```

#### Step 4: Compute χ² Statistic
Same as normal distribution test.

#### Step 5: Determine Degrees of Freedom
```
r = m - k - 1 = m - 0 - 1 = m - 1
```

Why k=0? Uniform distribution requires no parameter estimation.

#### Step 6: Compare and Decide
Same decision rule as before.

---

## Code Structure

### Data Structures

```go
// Stores coefficients for Goldstein approximation
type GoldsteinCoefficients struct {
    a, b, c float64
}

// Represents a single bin in the histogram
type Bin struct {
    Lower int   // Lower bound (inclusive)
    Upper int   // Upper bound (inclusive)
    Count int   // Number of observations
}

// Complete histogram
type Histogram struct {
    Bins       []Bin
    TotalCount int
}
```

### Main Program Flow

```
main()
  ↓
readInput()          → Get sample data and α
  ↓
createHistogram()    → Create and merge bins
  ↓
displayHistogram()   → Show frequency distribution
  ↓
testNormalDistribution()
  ├─ Estimate mean and std dev
  ├─ Calculate theoretical frequencies
  ├─ Compute χ² statistic
  ├─ Calculate χ²ₖᵣ using Goldstein
  └─ Compare and decide
  ↓
testUniformDistribution()
  ├─ Determine uniform range
  ├─ Calculate theoretical frequencies
  ├─ Compute χ² statistic
  ├─ Calculate χ²ₖᵣ using Goldstein
  └─ Compare and decide
```

---

## Mathematical Formulas

### 1. Chi-Squared Statistic

```
χ² = Σᵢ₌₁ᵐ [(nᵢ - nᵢᵀ)² / nᵢᵀ]
```

- m = number of bins
- nᵢ = observed frequency in bin i
- nᵢᵀ = expected frequency in bin i

### 2. Sample Mean

```
x̄ = (1/n) · Σᵢ₌₁ⁿ xᵢ
```

### 3. Sample Standard Deviation (Corrected)

```
s = √[(1/(n-1)) · Σᵢ₌₁ⁿ (xᵢ - x̄)²]
```

Note: Uses n-1 (Bessel's correction) for unbiased estimate.

### 4. Normal Distribution Probability

For bin [a, b]:
```
P(a < X < b) = Φ((b - μ)/σ) - Φ((a - μ)/σ)
```

Where Φ is the standard normal CDF:
```
Φ(z) = 0.5 · [1 + erf(z/√2)]
```

### 5. Uniform Distribution Probability

For uniform U[min, max] and bin covering k values:
```
P(bin) = k / (max - min + 1)
```

### 6. Degrees of Freedom

```
r = m - k - 1
```

- m = number of bins
- k = number of estimated parameters
- 1 = constraint (probabilities sum to 1)

### 7. Goldstein Approximation

```
χ²α,n = n · [Σᵢ₌₀⁶ n^(-i/2) · dⁱ · (aᵢ + bᵢ/n + cᵢ/n²)]³
```

Where d depends on α:
```
d = 2.0637 · (ln(1/(1-α)) - 0.16)^0.4274 - 1.5774     for 0.5 ≤ α ≤ 0.999
d = -2.0637 · (ln(1/α) - 0.16)^0.4274 + 1.5774         for 0.001 ≤ α < 0.5
```

---

## Usage Examples

### Example 1: Bell-Curve Data (Normal Distribution)

**Input:**
```
Sample: 9 9 9 9 9 9 10 10 10 10 10 10 10 10 10 10 10 10 10 10
        10 10 10 10 10 10 11 11 11 11 11 11 11 11 11 11 11 11
        11 11 11 11 11 11 11 11 11 11 11 11 11 11 11 11 11 11
        11 11 11 11 11 12 12 12 12 12 12 12 12 12 12 12 12 12 12 12
α: 0.05
```

**Histogram:**
```
[9]:  6 observations
[10]: 20 observations
[11]: 35 observations  ← peak
[12]: 15 observations
```

**Normal Distribution Test:**
```
Estimated mean: 10.7763
Estimated std dev: 0.8579

χ²ₑₘₚ = 1.0670
χ²ₖᵣ(0.05, 1) = 3.7665

Result: 1.0670 < 3.7665 → ACCEPTED ✓
```

**Interpretation:** The data follows a normal distribution. The small χ² value indicates observed frequencies closely match expected frequencies.

**Uniform Distribution Test:**
```
χ²ₑₘₚ = 23.2632
χ²ₖᵣ(0.05, 3) = 7.8137

Result: 23.2632 ≥ 7.8137 → REJECTED ✗
```

**Interpretation:** The data does NOT follow a uniform distribution. The large χ² value shows significant deviation from uniform.

---

### Example 2: Uniform Data

**Input:**
```
Sample: 0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20
        0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20
        0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20
α: 0.05
```

**Histogram (after merging):**
```
[0-1]:   6 observations
[2-3]:   6 observations
[4-5]:   6 observations
...
[18-20]: 9 observations  ← slightly more due to 3 values
```

**Uniform Distribution Test:**
```
χ²ₑₘₚ = 0.0000
χ²ₖᵣ(0.05, 9) = 16.9170

Result: 0.0000 < 16.9170 → ACCEPTED ✓
```

**Interpretation:** Perfect fit! Each bin has exactly the expected frequency.

**Normal Distribution Test:**
```
χ²ₑₘₚ = 13.7917
χ²ₖᵣ(0.05, 7) = 14.0651

Result: 13.7917 < 14.0651 → ACCEPTED ✓ (marginally)
```

**Interpretation:** Interestingly, uniform data can also pass the normal test with enough bins, though it's a marginal acceptance.

---

## Interpretation of Results

### Understanding χ² Values

| χ² Value | Interpretation |
|----------|----------------|
| Close to 0 | Perfect or near-perfect fit |
| Small (< critical) | Good fit, accept hypothesis |
| Close to critical | Marginal fit, borderline case |
| Large (> critical) | Poor fit, reject hypothesis |

### Understanding Significance Level (α)

| α | Meaning | Usage |
|---|---------|-------|
| 0.01 | Very strict (99% confidence) | When false positive is costly |
| 0.05 | Standard (95% confidence) | Most common in practice |
| 0.10 | Lenient (90% confidence) | Exploratory analysis |

**Lower α** → Harder to reject H₀ → More conservative
**Higher α** → Easier to reject H₀ → More strict

### Common Pitfalls

1. **Small sample size**: Need n ≥ 30 for reliable results
2. **Small bin frequencies**: Must have nᵢ ≥ 5 in each bin
3. **Too many/few bins**: Affects test power
4. **Wrong degrees of freedom**: Must account for estimated parameters

### When Test is Inconclusive

If χ²ₑₘₚ ≈ χ²ₖᵣ (very close):
- Test is borderline
- Consider:
  - Collecting more data
  - Using different α
  - Trying other goodness-of-fit tests (Kolmogorov-Smirnov, Anderson-Darling)

---

## Key Insights

### Why This Test Matters

1. **Validation**: Confirms assumptions about data distribution
2. **Model Selection**: Helps choose appropriate statistical models
3. **Quality Control**: Detects anomalies in processes
4. **Scientific Research**: Tests theoretical predictions

### Limitations

1. **Binning Dependency**: Results can depend on how bins are created
2. **Parameter Estimation**: Estimating parameters from same data reduces test power
3. **Discrete Data**: Originally designed for continuous distributions
4. **Multiple Testing**: Testing multiple distributions increases false positive risk

### Advantages

1. **Versatile**: Works with any distribution
2. **Easy to Understand**: Visual and intuitive
3. **Widely Used**: Standard statistical test
4. **Flexible**: Can test composite hypotheses

---

## Running the Program

```bash
cd c:\Users\rawr\GoProjects\emp\lb5
go run main.go
```

**Interactive Input:**
1. Enter sample data (space-separated integers 0-20)
2. Enter significance level α (e.g., 0.05)

**Output:**
- Histogram display
- Normal distribution test results
- Uniform distribution test results
- Acceptance/rejection decision with explanation

---

## References

- Pearson, K. (1900). "On the criterion that a given system of deviations from the probable in the case of a correlated system of variables is such that it can be reasonably supposed to have arisen from random sampling"
- Goldstein approximation for χ² quantiles
- Cornish-Fisher expansion (alternative to Goldstein)

---

## Summary

This implementation provides:
- ✓ Robust input validation
- ✓ Automatic bin merging for valid test conditions
- ✓ Goldstein approximation for critical values
- ✓ Complete normal distribution test
- ✓ Complete uniform distribution test
- ✓ Clear, interpretable output
- ✓ Educational examples and explanations

The Chi-Squared test is a powerful tool for validating distributional assumptions in statistical analysis!
