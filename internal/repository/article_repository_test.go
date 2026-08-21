package repository

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestArticleRepositorySearchBindsRelevanceArguments(t *testing.T) {
	parts := buildArticleSearchParts("Go")

	require.Len(t, parts.whereArgs, 5)
	require.Len(t, parts.relevanceArgs, 3)
	require.Contains(t, parts.relevanceSQL, "a.title LIKE ?")
	require.Equal(t, "%Go%", parts.whereArgs[0])
	require.Equal(t, "%Go%", parts.relevanceArgs[0])
}

func TestArticleRepositorySearchSplitsMultipleKeywordsInStableOrder(t *testing.T) {
	parts := buildArticleSearchParts("Go blog")

	require.Len(t, parts.whereArgs, 10)
	require.Len(t, parts.relevanceArgs, 6)
	require.Equal(t, "%Go%", parts.whereArgs[0])
	require.Equal(t, "%blog%", parts.whereArgs[5])
	require.Equal(t, "%Go%", parts.tagArgs[0])
	require.Equal(t, "%blog%", parts.tagArgs[1])
}

func TestArticleRepositorySearchEscapesLikeWildcards(t *testing.T) {
	parts := buildArticleSearchParts("100%_match")

	require.Equal(t, `%100\%\_match%`, parts.whereArgs[0])
}
