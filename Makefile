.PHONY: all build clean generate-service install-service uninstall-service

APP_NAME = snmp-wrapper
CURRENT_DIR = $(shell pwd)
LOGROTATE_CONF = /etc/logrotate.d/$(APP_NAME)
LOG_PATH = /var/log/$(APP_NAME)
SERVICE_FILE = /etc/systemd/system/$(APP_NAME).service

define LOGROTATE_CONF_CONTENT
$(LOG_PATH)/*.log {
    rotate 5
    weekly
    missingok
    notifempty
    compress
    create 644 root root
}
endef

export LOGROTATE_CONF_CONTENT


all: build

build:
	@echo "Building $(APP_NAME)..."
	go build -o $(CURRENT_DIR)/$(APP_NAME)

clean:
	@echo "Cleaning up..."
	-rm -f $(CURRENT_DIR)/$(APP_NAME)

generate-service:
	@read -p "Enter UDP host IP address to listen on for the SNMP server: " lHost; \
	read -p "Enter UDP port to listen on for the SNMP server: " lPort; \
	read -p "Enter Mikrotik API server URL (host:port): " rhost; \
	read -p "Enter Username for Mikrotik API authentication: " user; \
	read -p "Enter Password for Mikrotik API authentication: " password; \
	read -p "Enter the monitoring interval in seconds: " interval; \
	printf "%s\n" \
		"[Unit]" \
		"Description=SNMP Wrapper Service for MikroTik REST API" \
		"After=network.target" \
		"" \
		"[Service]" \
		"Type=simple" \
		"Restart=on-failure" \
		"RestartSec=30s" \
		"KillMode=process" \
		"Environment=LOGPATH_SNMP_WRAPPER=$(LOG_PATH)" \
		"ExecStart=/usr/local/bin/$(APP_NAME) -l $$lHost -p $$lPort -RHost $$rhost -U $$user -P $$password -i $$interval" \
		"ExecStop=bash -c \"ps aux | grep $(APP_NAME) | head -1 | awk {'print \$$2'} | xargs kill -9\"" \
		"ExecStop=/bin/kill $MAINPID" \
		"StandardError=append:$(LOG_PATH)/$(APP_NAME)_err.log" \
		"StandardOutput=append:$(LOG_PATH)/$(APP_NAME).log" \
		"" \
		"[Install]" \
		"WantedBy=multi-user.target" > $(SERVICE_FILE)

install-service: generate-service
	@echo "Installing $(APP_NAME) as a systemd service..."
	@echo "$$LOGROTATE_CONF_CONTENT" > $(LOGROTATE_CONF)
	@cp $(CURRENT_DIR)/$(APP_NAME) /usr/local/bin/$(APP_NAME)
	-mkdir -p $(LOG_PATH)
	@systemctl daemon-reload
	@systemctl enable $(APP_NAME)
	@systemctl start $(APP_NAME)
	@systemctl status $(APP_NAME)

uninstall-service:
	@echo "Uninstalling $(APP_NAME) service..."
	-systemctl stop $(APP_NAME)
	-systemctl disable $(APP_NAME)
	-rm -f $(SERVICE_FILE)
	@systemctl daemon-reload
	-rm -f /usr/local/bin/$(APP_NAME)
	-rm -rf $(LOG_PATH)
	@echo "$(APP_NAME) service uninstalled."

