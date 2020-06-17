import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { DetailComponent } from '../detail/detail.component';
import { MonacoEditorModule } from 'ng-monaco-editor';
import { LogsComponent } from './logs/logs.component';
import { MatTabsModule } from '@angular/material/tabs';

@NgModule({
  imports: [
    MonacoEditorModule,
    FormsModule,
    MatTabsModule,
  ],
  declarations: [
    DetailComponent,
    LogsComponent
  ],
})
export class DetailModule { }
