import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { ThemeService } from '../../theme/theme.service';
import { YamlViewerOptions } from '../../types/types';

@Component({
  selector: 'app-editor',
  templateUrl: 'editor.component.html'
})
export class EditorComponent implements OnInit {

  @Input() model: string;
  @Output() modelChange = new EventEmitter<string>();

  options: YamlViewerOptions;

  constructor(
    private themeService: ThemeService,
  ) { }

  ngOnInit() {
    this.options = {
      folding: true,
      minimap: { enabled: true },
      readOnly: true,
      language: 'yaml',
      contextmenu: false,
      scrollBeyondLastLine: false,
      lineNumbers: false as any,
      theme: this.isDarkTheme() ? 'vs-dark' : 'vs'
    };
    this.themeService.themeChanged.subscribe((newTheme: string) => {
      this.options = {
        ...this.options,
        theme: newTheme.includes('dark') ? 'vs-dark' : 'vs',
      };
    });
  }

  isDarkTheme(): boolean {
    return this.themeService.getTheme().includes('dark');
  }
}
