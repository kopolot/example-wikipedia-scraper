package repository

func handleNotFoundModel[T any](err error, result T) (T, error) {
	if err != nil {
		if err.Error() == "record not found" {
			return result, nil
		}
		return result, err
	}
	return result, nil
}

func handleNotFoundSlice[T any](err error, result []T) ([]T, error) {
	if err != nil {
		if err.Error() == "record not found" {
			return []T{}, nil
		}
		return nil, err
	}
	return result, nil
}
