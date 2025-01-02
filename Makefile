# Go parameters
GOCMD=go
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOMOD=$(GOCMD) mod

# Test flags
TEST_FLAGS=-v -race -count=1
TEST_TIMEOUT=30s

.PHONY: clean test test-cover tidy fmt

# 清理测试产物
clean:
	$(GOCLEAN)
	rm -f coverage.out

# 运行测试
test:
	$(GOTEST) $(TEST_FLAGS) -timeout $(TEST_TIMEOUT) ./...

# 运行测试并生成覆盖率报告
test-cover:
	$(GOTEST) $(TEST_FLAGS) -timeout $(TEST_TIMEOUT) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

# 整理依赖
tidy:
	$(GOMOD) tidy

# 格式化代码
fmt:
	$(GOCMD) fmt ./...