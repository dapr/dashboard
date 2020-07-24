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
  public showFiltered: boolean;
  public filterValue: string;

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
      this.logs = data;
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
