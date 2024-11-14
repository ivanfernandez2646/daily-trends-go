package domain

type FeedOption func(*Feed) error

func WithId(value string) FeedOption {
	return func(f *Feed) error {
		id, err := NewFeedId(value)
		if err != nil {
			return err
		}

		f.id = *id
		return nil
	}
}

func WithTitle(value string) FeedOption {
	return func(f *Feed) error {
		title, err := NewFeedTitle(value)
		if err != nil {
			return err
		}

		f.title = *title
		return nil
	}
}

func WithDescription(value string) FeedOption {
	return func(f *Feed) error {
		description := NewFeedDescription(value)
		f.description = description
		return nil
	}
}

func WithAuthor(value string) FeedOption {
	return func(f *Feed) error {
		author, err := NewFeedAuthor(value)
		if err != nil {
			return err
		}

		f.author = author
		return nil
	}
}

func WithSource(value string) FeedOption {
	return func(f *Feed) error {
		source, err := NewFeedSource(value)
		if err != nil {
			return err
		}

		f.source = source
		return nil
	}
}
