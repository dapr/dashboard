import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { DetailComponent } from '../detail/detail.component';
import { MonacoEditorModule } from 'ng-monaco-editor';
import { MatTabsModule } from '@angular/material/tabs';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatTableModule } from '@angular/material/table';
import { MatInputModule } from '@angular/material/input';
import { MatIconModule } from '@angular/material/icon';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { LogsComponent } from './logs/logs.component';
import { FilterPipe } from './logs/filter.pipe';
import { HighlightPipe } from './logs/highlight.pipe';
import { MatSelectModule } from '@angular/material/select';
import { ContainerPipe } from './logs/container.pipe';
import { SortPipe } from './logs/sort.pipe';

@NgModule({
  imports: [
    CommonModule,
    MonacoEditorModule,
    FormsModule,
    MatTabsModule,
    MatCardModule,
    MatButtonModule,
    MatTableModule,
    MatInputModule,
    MatIconModule,
    MatCheckboxModule,
    MatSelectModule,
  ],
  declarations: [
    DetailComponent,
    LogsComponent,
    FilterPipe,
    HighlightPipe,
    ContainerPipe,
    SortPipe,
  ],
})
export class DetailModule { }
