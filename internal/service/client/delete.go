package client

import "context"

func (s *service) Delete(ctx context.Context, uuid string) error {
	return s.clientRepository.Delete(ctx, uuid)
}
