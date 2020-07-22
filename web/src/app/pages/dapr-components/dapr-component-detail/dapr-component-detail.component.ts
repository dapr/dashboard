import { Component, OnInit } from '@angular/core';
import { DaprComponent } from 'src/app/types/types';
import { ComponentsService } from 'src/app/components/components.service';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-dapr-component-detail',
  templateUrl: './dapr-component-detail.component.html',
  styleUrls: ['./dapr-component-detail.component.scss']
})
export class DaprComponentDetailComponent implements OnInit {

  private name: string;
  public component: DaprComponent;
  public componentName: string;
  public componentJSON: string;
  public options: object;
  public loadedComponent;

  constructor(
    private route: ActivatedRoute,
    private componentsService: ComponentsService
  ) { }

  ngOnInit(): void {
    this.name = this.route.snapshot.params.name;
    this.componentsService.getComponent(this.name).subscribe((data: DaprComponent) => {
      console.log(data);
      this.component = data;
      this.componentName = data.metadata.name;
      this.componentJSON = JSON.stringify(data.spec.metadata, null, 2);
      this.loadedComponent = true;
    });
    this.options = {
      folding: true,
      minimap: { enabled: true },
      readOnly: false,
      language: 'yaml',
    };
  }
}
