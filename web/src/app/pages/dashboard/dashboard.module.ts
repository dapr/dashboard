import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { DashboardComponent } from './dashboard.component';
import { MatTableModule } from '@angular/material/table';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';

@NgModule({
  imports: [
    RouterModule,
    CommonModule,
    MatTableModule,
    MatCardModule,
    MatButtonModule,
  ],
  declarations: [
    DashboardComponent,
  ],
})
export class DashboardModule { }
