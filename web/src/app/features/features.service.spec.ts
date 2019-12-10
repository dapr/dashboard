import { TestBed } from '@angular/core/testing';

import { FeaturesService } from './features.service';

describe('FeaturesService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: FeaturesService = TestBed.get(FeaturesService);
    expect(service).toBeTruthy();
  });
});
