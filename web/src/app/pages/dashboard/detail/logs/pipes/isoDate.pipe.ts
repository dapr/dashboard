import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
    name: 'isoDate'
})
export class ISODatePipe implements PipeTransform {
    transform(nanoseconds: number): string {
        return new Date(nanoseconds / 1000000).toISOString();
    }
}
