#   Author: Pawananjani Kumar (pawananjani.kumar@goniyo.com)
#   CreatedAt: 28 Mar 2024
echo "generating coverage report..."
go tool cover -html=coverage.out
go tool cover -html=coverage.out -o coverage.html
