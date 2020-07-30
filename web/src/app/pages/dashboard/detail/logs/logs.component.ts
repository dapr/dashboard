import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { InstanceService } from 'src/app/instances/instance.service';
import { Log } from 'src/app/types/types';
import { MatSnackBar } from '@angular/material/snack-bar';

@Component({
  selector: 'app-logs',
  templateUrl: 'logs.component.html',
  styleUrls: ['logs.component.scss'],
})

export class LogsComponent implements OnInit {

  public logs: Log[];
  public id: string;
  public containers: string[];
  public showFiltered = false;
  public filterValue = '';
  public containerValue = '\[all containers\]';
  public dateOrder  = 'desc';
  public showTimestamps = false;

  constructor(
    private route: ActivatedRoute,
    private instances: InstanceService,
    private snackbar: MatSnackBar,
  ) { }

  ngOnInit(): void {
    this.id = this.route.snapshot.params.id;
    this.getLogs(false);
  }

  getLogs(refresh: boolean): void {
    this.instances.getLogs(this.id).subscribe((data: Log[]) => {
      const comparator = (a: Log, b: Log) => {
        return a.timestamp - b.timestamp;
      };
      this.logs = data.sort(comparator);
      this.containers = data.map(log => log.container).filter((value, index, self) => self.indexOf(value) === index);
      this.containers.unshift('\[all containers\]');
      if (refresh) {
        this.showSnackbar('Logs successfully refreshed');
      }
    });
  }

  showSnackbar(message: string): void {
    this.snackbar.open(message, '', {
      duration: 2000,
    });
  }
}
