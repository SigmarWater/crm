package customer

import "context"

func (s *service) Delete(ctx context.Context, uuid string) error {
	return s.customerRepository.Delete(ctx, uuid)
}
