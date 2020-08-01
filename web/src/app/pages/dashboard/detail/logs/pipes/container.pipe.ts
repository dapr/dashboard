import { Pipe, PipeTransform } from '@angular/core';
import { Log } from 'src/app/types/types';

@Pipe({
  name: 'container'
})
export class ContainerPipe implements PipeTransform {
  transform(items: Log[], container: string): Log[] {
    if (!items) { return []; }
    if (!container || container === '\[all containers\]') { return items; }
    return items.filter(item => {
      if (item && item.container) {
        return item.container === container;
      }
      return false;
    });
  }
}
