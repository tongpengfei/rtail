
BIN_PATH = ../bin

TARGET_RTAIL=rtail

TARGETS=$(TARGET_RTAIL)

all: clean $(TARGETS)

$(TARGET_RTAIL): main.go
	go build -mod vendor -o $@ $^

clean:
	rm -rf $(TARGETS)

install: all
	cp -rf $(TARGETS) $(BIN_PATH)/
