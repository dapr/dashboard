import { TestBed } from '@angular/core/testing';

import { ComponentsService } from './component.service';

describe('ComponentService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: ComponentsService = TestBed.get(ComponentsService);
    expect(service).toBeTruthy();
  });
});
