cp service/probe.service  /usr/lib/systemd/system/
systemctl daemon-reload
systemctl start probe
systemctl enable probe
