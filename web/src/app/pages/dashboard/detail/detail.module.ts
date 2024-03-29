import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { DetailComponent } from '../detail/detail.component';
import { MatTabsModule } from '@angular/material/tabs';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatTableModule } from '@angular/material/table';
import { MatInputModule } from '@angular/material/input';
import { MatIconModule } from '@angular/material/icon';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { LogsComponent } from './logs/logs.component';
import { FilterPipe } from './logs/pipes/filter.pipe';
import { HighlightPipe } from './logs/pipes/highlight.pipe';
import { MatSelectModule } from '@angular/material/select';
import { ContainerPipe } from './logs/pipes/container.pipe';
import { SortPipe } from './logs/pipes/sort.pipe';
import { TimeSincePipe } from './logs/pipes/timeSince.pipe';
import { MatNativeDateModule } from '@angular/material/core';
import { TimePipe } from './logs/pipes/time.pipe';
import { ISODatePipe } from './logs/pipes/isoDate.pipe';
import { SharedModule } from '../../../shared/shared.module';

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    MatTabsModule,
    MatCardModule,
    MatButtonModule,
    MatTableModule,
    MatInputModule,
    MatIconModule,
    MatCheckboxModule,
    MatSelectModule,
    MatDatepickerModule,
    MatNativeDateModule,
    SharedModule
  ],
  declarations: [
    DetailComponent,
    LogsComponent,
    FilterPipe,
    HighlightPipe,
    ContainerPipe,
    SortPipe,
    TimeSincePipe,
    TimePipe,
    ISODatePipe
  ],
})
export class DetailModule { }
