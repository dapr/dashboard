import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { DetailComponent } from '../detail/detail.component';
import { MonacoEditorModule } from 'ng-monaco-editor';
import { LogsComponent } from './logs/logs.component';
import { MatTabsModule } from '@angular/material/tabs';
import { MatCardModule } from '@angular/material/card';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { MatButtonModule } from '@angular/material/button';
import { MatTableModule } from '@angular/material/table';

@NgModule({
  imports: [
    CommonModule,
    MonacoEditorModule,
    FormsModule,
    MatTabsModule,
    MatCardModule,
    MatCheckboxModule,
    MatSnackBarModule,
    MatButtonModule,
    MatTableModule
  ],
  declarations: [
    DetailComponent,
    LogsComponent
  ],
})
export class DetailModule { }
