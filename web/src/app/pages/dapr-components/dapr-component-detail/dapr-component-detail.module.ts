import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { DaprComponentDetailComponent } from './dapr-component-detail.component';
import { MonacoEditorModule } from 'ng-monaco-editor';
import { MatTabsModule } from '@angular/material/tabs';
import { MatCardModule } from '@angular/material/card';
import { MatTableModule } from '@angular/material/table';

@NgModule({
  imports: [
    CommonModule,
    MonacoEditorModule,
    FormsModule,
    MatTabsModule,
    MatCardModule,
    MatTableModule,
  ],
  declarations: [
    DaprComponentDetailComponent
  ],
})
export class DaprComponentDetailModule { }
