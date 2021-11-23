
SUBDIR=src

BIN_PATH=./bin

define make_subdir
	@for i in $(SUBDIR); do\
		(make -C $$i $1) \
	done;
endef

all:
	$(shell test ! -d $(BIN_PATH) && mkdir $(BIN_PATH))
	$(call make_subdir)

clean:
	$(shell test -d $(BIN_PATH) && rm -rf $(BIN_PATH))
	$(call make_subdir, clean)

install:all
	$(call make_subdir, install)
