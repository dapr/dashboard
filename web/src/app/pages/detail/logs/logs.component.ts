import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { InstanceService } from '../../../instances/instance.service';
import { Log } from './log';
import { MatSnackBar } from '@angular/material/snack-bar';

@Component({
  selector: 'ngx-logs',
  templateUrl: './logs.component.html',
  styleUrls: ['./logs.component.scss'],
})

export class LogsComponent implements OnInit {

  public logs: Log[];
  public id: string;
  public info: boolean;
  public debug: boolean;
  public warning: boolean;
  public error: boolean;
  public fatal: boolean;

  constructor(
    private route: ActivatedRoute,
    private instances: InstanceService,
    private snackbar: MatSnackBar,
  ) { }

  ngOnInit() {
    this.id = this.route.snapshot.params.id;
    this.getLogs(false);
  }

  getLogs(refresh: boolean) {
    this.instances.getLogsArray(this.id).subscribe(data => {
      this.logs = data;
      if (refresh) {
        this.showSnackbar('Logs successfully refreshed');
      }
    });
  }

  isActive(level: string): boolean {
    if (level === "info") return this.info;
    if (level === "debug") return this.debug;
    if (level === "warning") return this.warning;
    if (level === "error") return this.error;
    if (level === "fatal") return this.fatal;
    return false;
  }

  showSnackbar(message: string) {
    this.snackbar.open(message, '', {
      duration: 2000,
    });
  }
}
